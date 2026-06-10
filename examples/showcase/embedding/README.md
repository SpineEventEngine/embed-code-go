# Embedding Examples

This folder is a runnable guide to embedding instructions. The positive
examples show supported features, and the negative examples show the failures a
user should expect when an instruction is malformed or stale.

Run the positive examples from the repository root:

```bash
go run ./main.go -mode check -config-path examples/showcase/embedding/embed-code.yml
```

## Feature Examples

### Simple Embedding

- [Whole file source](positive/whole-file-source.md) 
  shows how to embed the whole source file.
- [Paired instruction tag](positive/paired-instruction-tag.md)
  shows the preferred paired tag form used throughout the showcase.
- [Named source root](positive/named-source-root.md)
  shows how to use different configured source trees.

### Line, Range, And Glob Matching

- [Source line pattern](positive/source-line-pattern.md)
  embeds the first source line that matches a `line` pattern.
- [Start and end patterns](positive/start-end-pattern.md)
  embeds an inclusive source range selected by `start` and `end`.
- [Multi-line patterns](positive/multi-line-pattern.md)
  uses `\n` to match consecutive source lines.
- [Pattern escaping](positive/pattern-escaping.md)
  shows literal `*`, trailing `$`, leading `^`, and backslash-n text matches.

### Fragments

- [Named fragment](positive/named-fragment.md) 
  embeds a region wrapped with `#docfragment` and `#enddocfragment` markers.
- [Multi-part fragment separator](positive/multi-part-fragment-separator.md)
  joins repeated fragment parts with the configured separator.
- [Overlapping fragments](positive/overlapping-fragments.md) 
  shows fragment markers that share source lines.

### Rendered Content And Documents

- [Comment filtering](positive/comment-filtering.md) 
  keeps documentation comments while removing regular comments.
- [Markdown fence shielding](positive/markdown-fence-shielding.md) 
  shows that instruction-looking text inside ordinary code fences is inert.
- [HTML showcase](positive/html-showcase.html) 
  shows that HTML documents can be processed when the include patterns allow them.
- [Excluded document](positive/ignored-by-exclude.md) 
  shows that `doc-excludes` prevents a matching file from being processed.

## Negative Examples

The negative examples are intentionally broken and should fail in check mode.
Use them to recognize common diagnostics:

```bash
go run ./main.go -mode check -config-path examples/showcase/embedding/negative/processing-errors.yml
go run ./main.go -mode check -config-path examples/showcase/embedding/negative/stale.yml
```

The cases live in [negative/docs](negative/docs/).
