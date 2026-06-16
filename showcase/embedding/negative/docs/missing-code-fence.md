# Missing Code Fence

This scenario shows what happens when an active instruction has no owned code
fence.

## How It Fails

Every instruction must be followed immediately by a Markdown code fence. The
plain text line below is not a fence, so the parser reports the missing fence at
the instruction. In a real guide, add an opening and closing fence after the
instruction, even if the fence starts empty.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
This line is not a code fence.
