/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
