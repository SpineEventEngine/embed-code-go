package org.showcase;

// #docfragment "Greeter class"
public final class Greeting {
    private Greeting() {}

    // #docfragment "main()"
    public static void main(String[] args) {
        System.out.println(greeting("Ada"));
    }
    // #enddocfragment "main()"

    public static String greeting(String name) {
        return "Hello, " + name + "!";
    }
}
// #enddocfragment "Greeter class"
