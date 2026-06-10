package org.showcase;

/**
 * Creates public greetings.
 */
public interface CommentModes {
    /*
     * Internal implementation note.
     */
    String URL = "http://example.org/*not-comment*/";

    // Regular inline comment.
    String greet(String name); // trailing inline comment.
}
