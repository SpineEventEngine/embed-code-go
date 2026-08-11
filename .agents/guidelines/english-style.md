# English style

Grammar, punctuation, and spelling rules for English prose in source-code
comments and documentation. Layout and flow rules live in the
[writer skill](../skills/writer/SKILL.md); this page covers the language itself.

The catalog below is the single source of truth for three kinds of consumers:

- **Fixers** — the `proofread` skill edits existing text, applying a fix
  only when no leave-alone guard matches.
- **Reviewers** — documentation reviews flag catalog violations on changed lines.
- **Authors** — repository writing skills follow the same rules for new prose.

## Principles

1. **Fix errors, not taste.** Correct only what the catalog names as an
   error. Do not reword sentences that are already correct, and do not
   impose stylistic preferences the catalog does not state.
2. **Minimal edits.** Preserve the author's wording, meaning, and voice.
   The smallest change that fixes the error wins.
3. **When uncertain, leave it and report it.** A missed error is cheaper than
   a wrong fix. Record ambiguous occurrences for a human to decide rather than
   editing them automatically.
4. **Never touch code or machine-read text.** See the per-language list below.
5. **Consistency is per file.** Where correct alternatives exist (spelling
   dialect, serial comma), the unit of consistency is the file: align a
   clear outlier with the file's dominant convention; when no convention
   dominates, leave the text unchanged and report the mix.

## Where English prose lives

| Language                | Prose to check                                                         |
|-------------------------|------------------------------------------------------------------------|
| Kotlin / Java           | KDoc/Javadoc bodies and tag descriptions; block and line comments      |
| Protobuf                | `//` doc comments for types and fields, and file headers               |
| TypeScript / JavaScript | TSDoc/JSDoc bodies and tag descriptions; block and line comments       |
| Go                      | Doc comments immediately above declarations and other comments         |
| Markdown / AsciiDoc     | Body prose and headings                                                |

Within a doc comment, only the *description text* is prose. Tag machinery
is not: in `@param {string} userId The identifier of the user.`, only
"The identifier of the user." may be edited.

## Never edit

Some comment content is read by compilers and tools, or must stay verbatim for
other reasons. A purported fix there can change build behavior or break tooling.
Skip the following content:

**In every language**

- Code in any form: string literals, identifiers, fenced and indented code
  blocks, inline code spans, `{@code}` / `{@literal}` tags, `<pre>` /
  `<code>` HTML, `@sample` references, and documented commands.
- License and copyright headers — they are generated from the copyright
  profile and are outside project-owned prose.
- Quoted text (quotations stay verbatim), URLs, email addresses, version
  strings, and file paths.
- The machine-read TODO prefix; the description after it is ordinary prose.
- Editor and tool modelines (`// vim: …` and similar).

**Kotlin / Java**

- Inspection suppressions: `//noinspection …`.
- Javadoc/KDoc reference targets: the arguments of `{@link}`,
  `{@linkplain}`, and `@see`, and KDoc `[Symbol]` references.

**Protobuf**

- Lint directives such as `// buf:lint:ignore …`.

**TypeScript / JavaScript**

- Compiler directives: `// @ts-ignore`, `// @ts-expect-error`,
  `// @ts-nocheck`, and `/// <reference … />` triple-slash directives.
- Tool directives: `// eslint-disable*`, `// prettier-ignore`, coverage
  markers (`/* istanbul ignore … */`, `/* c8 ignore … */`,
  `/* v8 ignore … */`), webpack magic comments
  (`/* webpackChunkName: … */`), and `//# sourceMappingURL=…`.
- JSDoc machinery: tag names, `{Type}` braces, and parameter names.

**Go**

- Compiler and tool pragmas: any `//go:…` comment (`//go:build`,
  `//go:generate`, `//go:embed`, `//go:linkname`, …), the legacy
  `// +build` constraints, and `//nolint` directives. The absence of a
  space after `//` marks a pragma — never "fix" that spacing.
- cgo preambles: the comment block immediately above `import "C"` is
  compiled C code, not prose.
- The parsed `Deprecated:` marker word; its explanation is prose.
- The leading identifier of a doc comment: godoc requires the comment to
  start with the declared name (`// ServeHTTP handles …`). Never reword,
  lowercase, or move it; the catalog applies from the second word on.

**Markdown**

- YAML front matter, link-reference definitions, badge markup, and HTML
  comments used as directives (`<!-- markdownlint-disable -->`,
  `<!-- prettier-ignore -->`).

**AsciiDoc**

- Document attributes (`:attribute-name: value`) and `include::` macros.
- Conditional directives: `ifdef::`, `ifndef::`, `endif::`.
- Delimited blocks — listing (`----`), literal (`....`), and `[source]`
  blocks: the "fenced and indented code blocks" rule above covers Markdown
  fencing, not these AsciiDoc delimiters.

## Error catalog

Fixers group their reports by the topic headings below.

### Restrictive "which" vs. "that"

As a project writing rule, use "that" for a relative clause that *restricts
or identifies* its antecedent. Use "which", preceded by a comma, for a clause
that merely *adds information*. Russian *который* covers both restrictive and
non-restrictive senses, so direct transfer can blur this distinction.

| Before                         | After                                      |
|--------------------------------|--------------------------------------------|
| a plugin which forces versions | a plugin that forces versions              |
| the file, which is generated   | the file, which is generated — leave alone |

Replace "which" with "that", matching the original capitalization, only
when **none** of these guards applies:

- A comma (with optional whitespace) precedes "which" — non-restrictive.
- An opening parenthesis or a dash precedes it: `(which …`, `— which …`,
  `– which …`, `-- which …` — equally non-restrictive.
- A preposition precedes it: "in which", "of which", "with which",
  "by which", "to which", "at which", "from which", "on which",
  "for which", "into which", "upon which", "under which", "within which",
  "through which", "against which", "without which", and the like.
- Interrogative or determiner use, direct or embedded: "Which plugin…?",
  "decide which plugin", "depending on which mode", "no matter which", and
  "which of the following".
- Sentence-initial "Which".
- A hyphen abuts it — kebab-case names such as `which-fixer` are
  identifiers, not pronouns.
- Fixed phrases: the fused relative "that which" and the idiom "which is which".

### Articles

Russian has no articles, so missing or wrong articles are the most common
error class. Add the missing article inside full sentences; choose "a" vs.
"an" by the *sound* that follows, not the letter.

| Before                                   | After                                           |
|------------------------------------------|-------------------------------------------------|
| Returns value of given field.            | Returns the value of the given field.           |
| Throws exception if file does not exist. | Throws an exception if the file does not exist. |
| a HTTP request                           | an HTTP request                                 |
| an user, an unique key                   | a user, a unique key                            |
| a SDK, an URL                            | an SDK, a URL                                   |

Leave alone:

- Plural and uncountable nouns used generically: "returns metadata",
  "handles errors", "provides support for logging".
- Deliberately telegraphic contexts: headings, table cells, and changelog-style
  list fragments. Articles are enforced only in full sentences.
- Initialisms whose pronunciation varies: "a SQL query" (*sequel*) and
  "an SQL query" (*es-cue-el*) are both correct — align with the file's
  existing usage, otherwise leave.
- Bare identifiers used as names need no article, as in
  "Calls `close` after use".

### Subject–verb agreement

The verb agrees with the grammatical subject — watch the head noun of a
long subject, and "there is/are".

| Before                              | After                              |
|-------------------------------------|------------------------------------|
| The methods returns a copy.         | The methods return a copy.         |
| Each of the listeners are notified. | Each of the listeners is notified. |
| The list of errors are cleared.     | The list of errors is cleared.     |
| There is several options.           | There are several options.         |

Leave alone:

- "data" as singular or plural — both are accepted in technical writing;
  keep each file consistent.
- A backticked identifier as a subject names one object regardless of its
  grammatical number: "`options` holds the parsed flags." is correct.
- "a number of X are …" and "the number of X is …" have different subjects, so
  both forms are correct.

### Verb form in API summaries

The summary sentence of a function or method describes what the call does,
in the third-person singular: "Returns …", "Creates …", "Validates …".
An imperative opener is the error.

| Before                             | After                               |
|------------------------------------|-------------------------------------|
| `/** Return the current state. */` | `/** Returns the current state. */` |
| `// Copy copy the buffer.` (Go)    | `// Copy copies the buffer.` (Go)   |

Leave alone:

- Go's leading identifier (see **Never edit**); fix only its following verb.
- Imperative mood where it is the local convention: step-by-step instructions,
  tutorials, README commands, commit messages, and CLI help text.
- Type summaries written as noun phrases ("A thread-safe cache of …") or
  with "Represents …" — both correct.
- "This method returns …" is wordy but not an error; leave the text unchanged.
- A file that consistently uses another summary convention (imperative
  summaries are common in some JavaScript codebases) — report the file
  once instead of rewriting every summary.

### Prepositions

Fix only the pairs listed here; other verb–preposition pairings vary
legitimately and stay unchanged.

| Before                            | After                             |
|-----------------------------------|-----------------------------------|
| depends of                        | depends on                        |
| independent from                  | independent of                    |
| consists from                     | consists of                       |
| capable to handle                 | capable of handling               |
| waits the result                  | waits for the result              |
| listens the event                 | listens to the event              |
| in runtime, in compile time       | at runtime, at compile time       |
| on practice                       | in practice                       |
| on the screenshot, on the diagram | in the screenshot, in the diagram |
| typical for                       | typical of                        |
| access of the file                | access to the file                |

Leave alone:

- Transitive verbs that take a direct object: "awaits the result",
  "accesses the file", "enters the block", "discusses the design" are
  correct without a preposition.
- "listens for the event" — correct when the sense is awaiting a specific
  occurrence rather than subscribing; do not swap "for" to "to".
- Idioms: "in search of", "on the basis of".
- Dialect-linked variants: "different from" (universal), "different to"
  (British), "different than" (American, informal) — do not convert.

### Verb complementation

"Allow", "enable", and "permit" need an object before a to-infinitive;
without one, use a gerund or rephrase. "Recommend" and "suggest" take a
gerund, not a bare infinitive. (The Russian construction «позволяет сделать»
can transfer as the ungrammatical "allows to do".)

| Before                        | After                        |
|-------------------------------|------------------------------|
| allows to configure the build | allows configuring the build |
| enables to run tests          | enables running tests        |
| permits to access the field   | permits access to the field  |
| We recommend to use the DSL.  | We recommend using the DSL.  |
| suggest to add a test         | suggest adding a test        |
| It is worth to note           | It is worth noting           |

Leave alone:

- An object is present: "allows the caller to configure …" is correct.
- "allows for" + noun: "allows for customization" is correct.
- "It is recommended to use …" is correct because the infinitive complements
  the passive "recommended".
- "helps (to) do" — both forms are correct.
- "provides the possibility to …" — calque-flavored and wordy, but its
  grammaticality is disputed rather than settled: prefer "makes it
  possible to …" when writing new text, and report rather than auto-fix.

### Comparatives

A comparison carries its own function word: "greater **than** X", "equal
**to** X". Coordinating the two keeps both — "greater than or equal to X".
Dropping one or both yields "greater or equal X", a calque of the Russian
«больше или равно X». This direct transfer omits the English function words
because the Russian comparative instead governs the case of the noun.

| Before                               | After                                      |
|--------------------------------------|--------------------------------------------|
| the day is less or equal zero        | the day is less than or equal to zero      |
| a value greater or equal the limit   | a value greater than or equal to the limit |
| the size is equal or greater than 10 | the size is greater than or equal to 10    |

Leave alone:

- "greater than or equal to" and "less than or equal to" — already complete.
- Comparison operators inside code spans, such as `>=`, `<=`, and `a >= b`, are
  code tokens rather than prose.
- "no less than", "no more than", "at least", "at most" — idiomatic and
  complete without "or equal".
- "equals" used transitively: "the result equals zero" takes no "to".

Writing "then" for "than" is a separate error — see *Confusables*.

### Commas

Only the mechanical cases below are errors; most comma placement remains the
author's judgment and stays unchanged.

| Before                                    | After                                     |
|-------------------------------------------|-------------------------------------------|
| If the file is missing the build fails.   | If the file is missing, the build fails.  |
| The value is cached, it is computed once. | The value is cached; it is computed once. |
| The value that is returned, is cached.    | The value that is returned is cached.     |

- An **introductory subordinate clause** ("If …", "When …", "While …",
  "Unless …", "Because …", "Although …" followed by a subject and verb)
  takes a comma before the main clause. Short adverbial openers
  ("By default", "In this case") conventionally take one too, but adding
  it is optional — leave existing text alone.
- A **comma splice** joins two independent clauses with a bare comma.
  Prefer the minimal repair: a semicolon when the clauses are closely
  related, otherwise a period.
- **No comma between a subject and its verb**, however long the subject.
- The **serial (Oxford) comma** is a style choice, not an error: keep each
  list internally consistent, add one only when its absence is genuinely
  ambiguous, and never churn existing lists.

Leave alone: commas inside quoted text; the comma before "which" — the
which/that topic owns that decision; and "e.g." and "i.e." with or without a
following comma, which is dialect-linked.

### Hyphenated compound modifiers

Two words acting as one adjective *before* a noun are hyphenated.

| Before                    | After                     |
|---------------------------|---------------------------|
| read only mode            | read-only mode            |
| well known issue          | well-known issue          |
| case sensitive comparison | case-sensitive comparison |
| third party library       | third-party library       |
| long running task         | long-running task         |

Leave alone:

- Predicative position — after the noun and a linking verb ("the mode is
  read only") the hyphen is optional; do not edit either way.
- Adverbs in "-ly" never hyphenate: "fully qualified name" and
  "publicly available API" are correct.
- A number with a unit symbol ("a 5 GiB limit") takes no hyphen by convention.

### Confusables

Word pairs that spell-checkers miss because both are real words.

| Before                               | After                                 |
|--------------------------------------|---------------------------------------|
| The method returns it's result.      | The method returns its result.        |
| Let's you configure the build.       | Lets you configure the build.         |
| This value maybe null.               | This value may be null.               |
| more efficient then the default      | more efficient than the default       |
| Use this method to setup the server. | Use this method to set up the server. |
| Users can login with a token.        | Users can log in with a token.        |
| The server can not recover.          | The server cannot recover.            |

- The noun or adjective is one word; the verb is two: setup / set up,
  login / log in, backup / back up, shutdown / shut down,
  checkout / check out. ("Run the setup." vs. "Set up the server.")
- "cannot" is the standard form; "can not" is reserved for the rare
  emphatic "is able not to".
- "e.g." (for example) vs. "i.e." (that is), and "affect" vs. "effect",
  change meaning. Fix only when the context makes the intent unambiguous;
  otherwise report the occurrence.

### Punctuation and spacing

| Before                    | After                    |
|---------------------------|--------------------------|
| `/** Returns the ID */`   | `/** Returns the ID. */` |
| the value.··The next step | the value.·The next step |
| the value ; the next      | the value; the next      |
| ends here..               | ends here.               |
| ( the default )           | (the default)            |

(`·` marks a space character.)

- A doc-comment summary sentence ends with a period — Javadoc and Dokka
  use it to delimit the summary.
- One space between sentences; no space *before* `.`, `,`, `;`, `:`, `?`,
  `!`; no space just inside parentheses or brackets in prose.
- No duplicated terminal punctuation ("..", "!!", "??"). The three-dot
  ellipsis "..." and the character "…" are legitimate; only a stray
  two-dot sequence is a typo.
- Sentences start with a capital letter — but **never change the case of
  a code identifier** to achieve this. Reword instead: "`timeout` limits
  the wait." → "The `timeout` value limits the wait."
- Missing apostrophes in contractions are spelling errors, such as "dont" → "don't".
- List items may be fragments without terminal periods — punctuate each
  list consistently and leave the chosen style alone.

### Spelling and dialect

Fix genuine misspellings outright:

| Before                      | After                        |
|-----------------------------|------------------------------|
| recieve, occured, seperate  | receive, occurred, separate  |
| existance, paramter, lenght | existence, parameter, length |
| successfull, usefull        | successful, useful           |
| compatable, preferrable     | compatible, preferable       |

For words where American and British English differ, **the unit of consistency is the file**:

- A clear outlier aligns with the file's dominant dialect: one
  "behaviour" in a file that otherwise writes "behavior" becomes
  "behavior" — and vice versa.
- When no dialect clearly dominates (the file is split), leave the text
  unchanged and report the mix.
- Judge dominance by strong markers: -or/-our (color/colour), -er/-re
  (center/centre), the licence/license noun, gray/grey. Treat -ise/-ize
  as a weak signal — British Oxford spelling legitimately uses "-ize", so
  "-ize" beside "behaviour" is not an inconsistency.
- Never change dialect in identifiers or code (a `Color` API stays
  `Color`), in quoted text, in proper nouns, or in prose that names the
  concept an identifier spells: a field called `colour` is described as
  "the colour", whatever the file's dialect.

There is no organization-wide dialect; if one is adopted later, the
conversion becomes a separate, mechanical sweep.

