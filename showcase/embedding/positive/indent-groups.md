# Fragment Indentation Groups

Use `indent-group` when a multi-part fragment combines independently indented
source regions. Each group receives its own common-indentation baseline, while
all partitions in that group preserve their indentation relative to one another.

## How It Works

Add `indent-group="name"` to an opening `#docfragment` marker. The group name
must be a non-empty quoted string. The matching `#enddocfragment` marker does
not repeat the attribute.

Partitions without `indent-group` belong to one shared default group. Existing
fragment markers therefore keep the standard behavior of normalizing common
indentation across the complete fragment.

In [GroupedIndent.java](../../code/java/org/showcase/GroupedIndent.java), the
import belongs to the `imports` group. The other three partitions use the
default group, so the top-level import does not affect their shared baseline:

```java
// #docfragment "Grouped example" indent-group="imports"
import java.util.List;
// #enddocfragment "Grouped example"

public static void print(List<String> values) {
    // #docfragment "Grouped example"
    var first = values.get(0);
        var normalized = first.trim();
    // #enddocfragment "Grouped example"

    // Two more "Grouped example" partitions use the default group.
}
```

## Embedding Instruction

The embedding instruction still selects only the fragment name. Indentation
groups are source-marker metadata and require no instruction attribute.

<embed-code file="$java/org/showcase/GroupedIndent.java" fragment="Grouped example"></embed-code>
```java
import java.util.List;
// ...
var first = values.get(0);
    var normalized = first.trim();
// ...
var second = values.get(1);
// ...
System.out.println(normalized + second);
```
