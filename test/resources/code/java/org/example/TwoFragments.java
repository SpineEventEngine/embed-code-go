package org.example;

// #docfragment "Main"
public class TwoFragments {
    // #enddocfragment "Main"
    private TwoFragments() {}

    // #docfragment "Main"
    public static void main(String[] args) {
        // #enddocfragment "Main"
        System.out.println("This is just a log message.");
        // #docfragment "Main"
        System.out.println(helperMethod());

    }
    // #enddocfragment "Main"

    // #docfragment "Hello"
    public static void hello(String[] args) {
        // #enddocfragment "Hello"
        var unseenText = "Unseen Text";
        System.out.println(unseenText);
        // #docfragment "Hello"
        var coolText = "Cool Text";
        System.out.println(coolText);
    }
    // #enddocfragment "Hello"

    private static String helperMethod() {
        return "42";
    }
// #docfragment "Main"
}
// #enddocfragment "Main"
