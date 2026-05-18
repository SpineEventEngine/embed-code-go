package org.example;

/**
 * Documents the public API.
 */
public interface Comments {
    /*
     * The block comment.
     */
    String marker = "http://example.org/*not-comment*/";

    // Full-line inline comment.
    String create(String name); // end-of-line inline comment.
}
