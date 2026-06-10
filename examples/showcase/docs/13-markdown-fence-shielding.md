# Instructions Inside Markdown Fences

Instruction-looking text inside an ordinary Markdown code fence is preserved as
documentation content. It is not executed because the parser tracks Markdown
fence state before looking for instructions.

````markdown
<embed-code file="$java/org/showcase/DoesNotRun.java"></embed-code>
```go
```
````
