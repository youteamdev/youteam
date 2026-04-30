set positional-arguments

dev *args:
    #!/usr/bin/env bash
    set -euo pipefail

    mkdir -p dist

    semver_pattern='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(\.(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?(\+([0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*))?$'
    version="0.0.0"
    tag="$(git describe --tags --exact-match HEAD 2>/dev/null || true)"
    normalized_tag="${tag#v}"
    if [[ "$normalized_tag" =~ $semver_pattern ]]; then
        version="$tag"
    fi
    commit="$(git rev-parse --short=7 HEAD 2>/dev/null || printf dev)"

    go build -ldflags "-X main.version=${version} -X main.commit=${commit}" -o dist/youteam ./cmd/youteam
    ./dist/youteam "$@"

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
