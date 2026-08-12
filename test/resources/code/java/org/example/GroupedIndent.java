package org.example;

// #docfragment "Example" indent-group="imports"
import java.util.List;
// #enddocfragment "Example"

public final class GroupedIndent {
    private GroupedIndent() {}

    // #docfragment "Example" indent-group="imports"
    static final String LABEL = "value";
    // #enddocfragment "Example"

    static void render(List<String> values) {
        // #docfragment "Example"
        var first = values.get(0);
            var nested = first.trim();
        // #enddocfragment "Example"

        // #docfragment "Example"
        var second = values.get(1);
        // #enddocfragment "Example"

        // #docfragment "Example" indent-group="output"
        System.out.println(nested + second);
        // #enddocfragment "Example"
    }
}
