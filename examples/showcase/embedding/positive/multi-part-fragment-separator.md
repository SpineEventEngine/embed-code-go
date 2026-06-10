# Multi-Part Fragment Separator

One named fragment can be split across several source regions.

## How It Works

[../../code/java/org/showcase/MultiPartWorkflow.java](../../code/java/org/showcase/MultiPartWorkflow.java)
opens and closes the `Workflow` fragment multiple times. Embed mode collects
each part in source order and joins the parts with the `separator` configured in
[../embed-code.yml](../embed-code.yml). This is useful when a guide needs the
shape of a class but wants to hide internal lines between selected regions.

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
