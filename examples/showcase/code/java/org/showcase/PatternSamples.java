package org.showcase;

class PatternSamples {

    private static final String ESCAPED_NEWLINE = "\n";

    @Scenario
    @Name("adds two numbers")
    void addsTwoNumbers() {
        int total = 1 + 1;

        assertEquals(2, total);
    }

    @Scenario
    @Name("subtracts two numbers")
    void subtractsTwoNumbers() {
        int total = 2 - 1;

        assertEquals(1, total);
    }
}
