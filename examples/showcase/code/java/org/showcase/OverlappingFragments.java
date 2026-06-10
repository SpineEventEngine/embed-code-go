package org.showcase;

// #docfragment "Class wrapper", "Greeting method"
public final class OverlappingFragments {
    // #enddocfragment "Class wrapper", "Greeting method"
    private OverlappingFragments() {}

    // #docfragment "Class wrapper"
    public static void boot() {
        // #enddocfragment "Class wrapper"
        System.out.println("Boot details are hidden.");
        // #docfragment "Class wrapper"
        System.out.println("Boot complete");
    }
    // #enddocfragment "Class wrapper"

    // #docfragment "Greeting method"
    public static String greeting(String name) {
        // #enddocfragment "Greeting method"
        var normalized = name.trim();
        // #docfragment "Greeting method"
        return "Hello, " + normalized + "!";
    }
    // #enddocfragment "Greeting method"

// #docfragment "Class wrapper", "Greeting method"
}
// #enddocfragment "Class wrapper", "Greeting method"
