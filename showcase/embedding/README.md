# Embedding Instructions

This folder is a runnable guide to `<embed-code>` instructions. The positive
examples show supported features, and the negative examples show failures users
should expect when an instruction is malformed or stale.

Run the positive examples from the repository root:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

## Feature Examples

### Instruction Shape

- [instruction-tag.md](positive/instruction-tag.md)
  explains managed code fences, supported tag forms, attributes, and source
  selection rules.
- [whole-file-source.md](positive/whole-file-source.md)
  embeds a complete source file.
- [named-source-root.md](positive/named-source-root.md)
  selects one configured source tree with `$name/`.

### Line, Range, And Glob Matching

- [source-line-pattern.md](positive/source-line-pattern.md)
  embeds the first source line that matches a `line` pattern.
- [start-end-pattern.md](positive/start-end-pattern.md)
  embeds an inclusive source range selected by `start` and `end`.
- [multi-line-pattern.md](positive/multi-line-pattern.md)
  uses `\n` to match consecutive source lines.
- [pattern-escaping.md](positive/pattern-escaping.md)
  shows how to match literal glob characters.

### Fragments

- [named-fragment.md](positive/named-fragment.md)
  embeds a region wrapped with `#docfragment` and `#enddocfragment` markers.
- [multi-part-fragment-separator.md](positive/multi-part-fragment-separator.md)
  joins repeated fragment parts with the configured separator.
- [overlapping-fragments.md](positive/overlapping-fragments.md)
  shows fragment markers that share source lines.

### Rendered Content And Documents

- [comment-filtering.md](positive/comment-filtering.md)
  compares `comments` modes and lists language support.
- [markdown-fence-shielding.md](positive/markdown-fence-shielding.md)
  shows that instruction-looking text inside ordinary code fences is inert.
- [html-showcase.html](positive/html-showcase.html)
  shows that HTML documents can be processed when include patterns allow them.

## Negative Examples

The negative examples are intentionally broken and should fail. Use them to
recognize common diagnostics:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/negative/processing-errors.yml
go run ./main.go -mode=check -config-path=showcase/embedding/negative/stale.yml
```

The cases live in [negative/docs](negative/docs).
