set positional-arguments

dev *args:
    @mkdir -p dist
    @go build -o dist/youteam ./cmd/youteam
    @./dist/youteam "$@"
