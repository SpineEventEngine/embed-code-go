# Multi-Part Fragment Separator

Fragments can be multipart, if fragment with the same name 
is started and ended in the file several times.

## How It Works

[../../code/java/org/showcase/MultiPartWorkflow.java](../../code/java/org/showcase/MultiPartWorkflow.java)
opens and closes the `Workflow` fragment multiple times. 
Embed mode collects each part in source order and joins the parts 
with the `separator` configured in [../embed-code.yml](../embed-code.yml). 

Separator is `...` by default.

<embed-code file="$java/org/showcase/MultiPartWorkflow.java" fragment="Workflow"></embed-code>
```java
public final class MultiPartWorkflow {
    // ...
    public static void start() {
        // ...
        System.out.println("Start workflow");
    }
// ...
}
```
