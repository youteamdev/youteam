set positional-arguments

dev *args:
    @mkdir -p dist
    @go build -o dist/youteam ./cmd/youteam
    @./dist/youteam "$@"

test:
    #!/usr/bin/env bash
    set -euo pipefail

    coverage_file="$(mktemp)"
    trap 'rm -f "$coverage_file"' EXIT

    go test ./... -covermode=atomic -coverprofile="$coverage_file"

    coverage="$(go tool cover -func="$coverage_file" | awk '/^total:/ { gsub(/%/, "", $3); print $3 }')"
    if [[ -z "$coverage" ]]; then
        echo "failed to calculate total coverage" >&2
        exit 1
    fi

    awk -v coverage="$coverage" 'BEGIN {
        if (coverage + 0 < 90) {
            printf "coverage %.1f%% is below 90%%\n", coverage
            exit 1
        }
        printf "coverage %.1f%% meets 90%% threshold\n", coverage
    }'
