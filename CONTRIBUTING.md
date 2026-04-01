# Contributing to Go-Charts

## Prerequisites

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/)
- go-gui (sibling directory at `../go-gui`)

## Build and Test

```bash
go build ./...             # build all packages
go test ./...              # run all tests
go vet ./...               # static analysis
golangci-lint run ./...    # full lint
```

## Coding Conventions

- **No variable shadowing.** Use `=` to reassign existing variables, not `:=`.
- **Clean lint and format.** All code must pass `golangci-lint run ./...`
  and `gofmt` with zero issues before committing.
- All chart types follow the `*Cfg` struct pattern from go-gui.
- Charts implement `gui.View` interface.

## Submitting Changes

1. Fork the repository and create a feature branch.
2. Make focused, single-purpose commits.
3. Add or update tests for any changed behavior.
4. Run the full check suite before pushing.
5. Open a pull request against `main`.

## License

Contributions are accepted under the
[PolyForm Noncommercial License 1.0.0](LICENSE).
