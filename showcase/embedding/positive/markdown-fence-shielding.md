# Instructions Inside Markdown Fences

Use a normal Markdown fence when documentation needs to show an embedding
instruction as plain text. This lets guides explain the syntax without causing
the example instruction to run.

## How It Works

The parser tracks ordinary code fences before it looks for active
`<embed-code>` instructions. Instruction-looking text inside a fence is
therefore preserved as documentation content, not treated as a real instruction.

An active instruction must appear outside a fence and must be followed by its
own managed code fence. The nested fence in this example is only part of the
displayed Markdown snippet.

## Shielded Example

````markdown
<embed-code file="$java/org/showcase/DoesNotRun.java"></embed-code>
```go
```
````
