# Instructions Inside Markdown Fences

Documentation sometimes needs to show an instruction as plain text.

## How It Works

The instruction-looking text below is inside an ordinary Markdown fence, so it
is preserved as documentation content. The parser tracks code-fence state before
looking for active instructions, which prevents examples from accidentally
running while they are being explained.

````markdown
<embed-code file="$java/org/showcase/DoesNotRun.java"></embed-code>
```go
```
````
