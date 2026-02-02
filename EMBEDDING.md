# Setting up documentation and code files

Embed Code handles a custom tag `<embed-code>` that allows you to embed code
samples from your code files into your documentation. 

There are two ways to specify the code fragment to embed.

**Option 1: A named fragment**
```
<embed-code file="path/to/file" fragment="Fragment Name"></embed-code> (I)
```

**Option 2: Regular expressions**
```
<embed-code file="path/to/file" start="first?line*glob" end="last?line*glob"></embed-code> (II)
```

The instruction must always be followed by a code fence (opening and closing three backticks):
<pre>
<embed-code ...></embed-code>
```java
```
</pre>

The content of the code fence does not matter — the command will overwrite it automatically.

Note that the code fence may specify the syntax in which the code will be highlighted.

This is true even when embedding into HTML.

## Named fragments (I)

## Markup fragments

You can mark up the code file to select named fragments like this:

```java
public final class String
    implements java.io.Serializable, Comparable<String>, CharSequence {
    
    // #docfragment "Constructor"
    public String() {
        this.value = new char[0];
    }
    // #enddocfragment "Constructor"
}
```

The `#docfragment` and `#enddocfragment` tags won't be copied into the resulting code fragment.

## Add embedding instructions

To add a new code sample, add the following construct to the Markdown file:

<pre>
&lt;embed-code file=&quot;java/lang/String.java&quot;
             fragment=&quot;Constructor&quot;&gt;&lt;/embed-code&gt;
```java
```   
</pre>

The `file` attribute specifies the path to the code file relative to the code root, specified in
the configuration. The `fragment` attribute specifies the name of the code fragment to embed. Omit
this attribute to embed the whole file.

You may use any name for your fragments, just omit double quotes (`"`) and symbols forbidden in XML.

## Pattern fragments (II)

Alternatively, the `<embed-code>` tag may have the following form:
<pre>
&lt;embed-code file=&quot;java/lang/String.java&quot;
             start=&quot;*class Hello*&quot;
             end=&quot;}*&quot;&gt;&lt;/embed-code&gt;
```java
```   
</pre>

In this case, the fragment is specified by a pair of glob-style patterns. The patterns match
the first and the last lines of the desired code fragment. Any of the patterns may be skipped.
In such a case, the fragment starts at the beginning or ends at the end of the code file.

The pattern syntax supports an extended glob syntax:
- `?` — one arbitrary symbol;
- `*` — zero, one, or many arbitrary symbols;
- `[set]` — one symbol from the given set (equivalent to `[set]` in regular expressions);
- `^` at the start of the pattern to signify the start of the line;
- `$` at the end of the pattern to signify the end of the line.

Note that the `*` symbols at the start and in the end of the pattern are implied. Use `^` and `$` to
mark that the pattern should not assume `*` at the start/end.

Note that `^` and `$` work as special characters in their respective positions but not in the middle
of the pattern. To match the literal `^` symbol at the start of the line, prepend it with another
`^`. Similarly, to match a literal `$` at the end of the line, append it with another `$`.

