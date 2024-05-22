# General
This is the implementation of embed-code written in Go. For more information, you can check the [README.md](https://github.com/SpineEventEngine/embed-code/blob/master/README.md).

Embed-code-go provides a way to embed code snippets into markdown files. This allows developers to easily include code examples within their documentation. 

Key features:
- extracts code fragments from source files and embeds them into documentation files;
- verifies that embedded code samples are up-to-date with the source code;
- supports flexible configuration via command line arguments or a YAML configuration file;
- allows embedding specific code fragments or code between start and end line patterns.

# Running

Embed-code-go can be run in three modes:
1. Embedding — in this mode, all documentation files are scanned for `<embed-code>` tags, and the corresponding embeddings are performed. Embedding results are written to the given documentation files.
2. Checking for being up-to-date — in this mode, all documentation files are checked for being up-to-date. All the embeddings under `<embed-code>` tags are compared with the corresponding code fragments. If the code is completely equals to the corresponding embeddings, the documentation files are up-to-date. If the check fails, an error message with the corresponding files is provided.
3. Analyzing — in this mode, all embeddings are checked for having corresponding code fragments. All information about problems found is written to the `build/analytics/problem-files.txt` file.

The mode is selected by the mandatory `mode` argument. If it is set to `check`, the checking for up-to-date is performed. If it is set to `embed`, the embedding is performed. If it is set to `analyze`, the analyzing is performed.

The tool can be executed as a binary file or as a Go file. In the latter case, the user must have Go [installed](#installation). The binary files are stored in the `./bin` directory.

The code and documentation files must be prepared for embedding. The instructions are provided in the [Setting up documentation and code files](#setting-up-documentation-and-code-files) section.

### Running binary executable
To run the `embed_code` binary executable, the following command can be used:
```
./<binary_executable_name> [arguments]
```
The binaries are located in the `./bin` directory. 

### Running Go file

#### Go version

Make sure you have Go [installed](#installation). Our version is `1.22.1`.

#### Running

To run the `embed_code.go` file, the following command can be used:
```
go run ./embed_code.go [arguments]
```

### Arguments
The available arguments are:
  * `-mode`: mandatory, `check` to checking for code embeddings to be up-to-date; `embed` to start the embedding process; `analyze` to run analyzing;
  * `-code_root`: optional, path to the root directory containing code files;
  * `-docs_root`: optional, path to the root directory containing documentation files;
  * `-config_file_path`: optional, path to a YAML configuration file that contains the code_root and docs_root fields;
  * `-code_includes`: optional, a comma-separated string of glob patterns for code files to include. For example: `"**/*.java,**/*.gradle"`. Default value is `"**/*.*"`;
  * `-doc_includes`: optional, a comma-separated string of glob patterns for docs files to include. For example: `"docs/**/*.md,guides/*.html"`. Default value is `"**/*.md,**/*.html"`;
  * `-fragments_dir`: optional, a path to a directory with code fragments. Default value is `./build/fragments`;
  * `-separator`: optional, a string which is used as a separator between code fragments. Default value is `...`.
 
Even though the `code_root`, `docs_root`, and `config_file_path` arguments are optional, Embed-code still requires the root directories for code and documentation to be set. This can be done in one of two ways:

1. Provide the `code_root` and `docs_root` arguments, in this case the roots are read directly from the provided paths.
2. Provide the `config_file_path` argument, in this case the roots are read from the given configuration file.

If neither of these options is provided, the embedding process will fail. If both options are set, the embedding will also fail.

### Configuration file

Optional settings can be defined in the configuration file. The file is a YAML file with the following structure:

```yaml
code_root: path/to/code/root
docs_root: path/to/docs/root
code_includes: "**/*.java,**/*.gradle"
```

The available arguments for the config file are:
  * `code_root`: mandatory;
  * `docs_root`: mandatory;
  * `config_file_path`: optional;
  * `code_includes`: optional;
  * `doc_includes`: optional;
  * `fragments_dir`: optional;
  * `separator`: optional.

These settings have the same role as the command-line arguments.

## Setting up documentation and code files

Synopsis:
```
<embed-code file="path/to/file" fragment="Fragment Name"></embed-code> (I)

OR

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

#### Named fragments (I)

##### Markup fragments

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

##### Add embedding instructions 

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

#### Pattern fragments (II)

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

# Installation

* Go to https://go.dev/doc/install. Our Go version is `1.22.1`, which can be checked in the [go.mod](./go.mod) file
* Make sure your Go installed successfully with the command
    ```bash
    go version
    ```

# Compilation
The pre-compiled binary executables are stored in the `./bin` directory. However, it is also possible to compile the file manually.
To compile the file, ensure that the Go is [installed](#installation).

Open terminal and navigate to the directory where `embed_code.go` is located. Then, use the following command to compile the file: 
```
go build embed_code.go
```

There may be issues when running `go build` outside of the directory containing `embed_code.go`, even if the path is specified correctly.

This command will create an executable file named `embed_code` (or `embed_code.exe` on Windows) in the same directory.
For further information, please refer to the [docs](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies).
