package org.showcase;

// #docfragment "Workflow"
public final class MultiPartWorkflow {
    // #enddocfragment "Workflow"
    private MultiPartWorkflow() {}

    // #docfragment "Workflow"
    public static void start() {
        // #enddocfragment "Workflow"
        System.out.println("Internal setup is hidden.");
        // #docfragment "Workflow"
        System.out.println("Start workflow");
    }
    // #enddocfragment "Workflow"

    public static void finish() {
        System.out.println("Finish workflow");
    }

// #docfragment "Workflow"
}
// #enddocfragment "Workflow"
