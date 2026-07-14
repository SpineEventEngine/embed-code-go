# Showcase End-to-End Tests

The showcase can also run as an end-to-end test suite. These steps are intended
for Embed Code application developers and are not part of the application usage
guide.

Run commands from the repository root.

The end-to-end suite checks positive examples, expected failures, and all
configuration shapes:

```bash
go test -tags showcase ./showcase
```

Check the quick start directly:

```bash
cd showcase/quick-start
go run ../../main.go -mode=check -config-path=config.yml
```

Check the positive embedding examples directly:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

Check the configuration examples directly:

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/single-source.yml
go run ./main.go -mode=check -config-path=showcase/configuration/named-sources.yml
go run ./main.go -mode=check -config-path=showcase/configuration/include-exclude.yml
go run ./main.go -mode=check -config-path=showcase/configuration/multiple-embeddings.yml
```

The negative examples are intentionally broken, so these commands should fail:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/negative/processing-errors.yml
go run ./main.go -mode=check -config-path=showcase/embedding/negative/stale.yml
```
