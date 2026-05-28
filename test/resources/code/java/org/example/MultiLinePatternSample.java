package org.example;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

class MultiLinePatternSample {

    private static final String MY_STRING = "\n";

    @Test
    @DisplayName("adds two values")
    void addsTwoValues() {
        int value = 1 + 1;

        assertEquals(2, value);
    }

    @Test
    @DisplayName("subtracts two values")
    void subtractsTwoValues() {
        int value = 2 - 1;

        assertEquals(1, value);
    }
}
