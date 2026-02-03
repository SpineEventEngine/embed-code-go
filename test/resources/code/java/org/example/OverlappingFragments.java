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
package org.example;

// #docfragment "Main", "Hello"
public class OverlappingFragments {
    // #enddocfragment "Main", "Hello"
    private OverlappingFragments() {}

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
// #docfragment "Main", "Hello"
}
// #enddocfragment "Main", "Hello"
