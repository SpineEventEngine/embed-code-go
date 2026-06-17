# Technical Audit & Improvement Plan — `embed-code-go`

| | |
|---|---|
| **Repository** | `SpineEventEngine/embed-code-go` |
| **Audit date** | 2026-06-17 |
| **Scope** | Full-repository audit. Analysis only — no source code was modified. |
| **Method** | Every source file read; build, `go vet`, and the full test suite run; coverage measured; the headline correctness bug reproduced from a verbatim copy of the affected functions. |
| **Citations** | Findings cite `file:line`. Each is labeled **[Fact]** (verifiable in code / reproduced) or **[Judgment]** (design opinion). |
| **Overall grade** | **B+ (Good)** |

---

## Table of contents
1. [Decisions taken (owner sign-off)](#1-decisions-taken-owner-sign-off)
2. [Executive summary](#2-executive-summary)
3. [Repo map](#3-repo-map)
4. [Audit report](#4-audit-report)
5. [Improvement strategy](#5-improvement-strategy)
6. [Task plan](#6-task-plan)
7. [Appendix — evidence](#7-appendix--evidence)

---

## 1. Decisions taken (owner sign-off)

These were open questions in the original audit; the product owner has since resolved them. They are now constraints on the plan.

| Topic | Decision |
|---|---|
| **Binary distribution** | Move from committed `bin/` to **GitHub Releases**. |
| **Repository history** | **Purge the entire `bin/` history** to reclaim space (accepts a destructive, repo-wide rewrite). |
| **Audience** | The tool **will be consumed by external teams** (not internal-only). |
| **License** | **Apache-2.0.** (The current per-file "All rights reserved" headers do not grant external use, so a root `LICENSE` is now required.) |
| **Comment-filter correctness** | Reliable behavior is **desired for these languages, in priority order:** Kotlin → Java → TypeScript → JavaScript → C# → C++ → Python → Go → Visual Basic. |
| **`check` output contract** | **No** downstream pipeline depends on `check` stdout/exit format — a refactor may change it freely. |

**One coordination item remains:** the history purge (task **M2-1b**) requires temporarily lifting `master` branch protection and a maintainer-performed force-push. It is *not* normal feature work and should be scheduled deliberately.

---

## 2. Executive summary

`embed-code-go` is a mature, well-engineered Go CLI that injects source snippets into Markdown documentation (a Go rewrite of a Ruby `embed-code` utility) for Hugo-based docs. The architecture is clean and layered — a State-pattern line parser feeding a fragment resolver and per-language comment filter — and the project enforces an unusually strict 40-linter `golangci-lint` configuration in CI across Windows/macOS/Linux. It ships typed errors with clickable `file://` diagnostics and carries **75.9% integrated test coverage** with tests that assert real behavior rather than mere execution. The engineering bar is clearly above average.

Three issues hold it back from an A. **Top 3 risks:** (1) a **confirmed runtime panic** in `indent.CutIndent` (`indent/indent.go:65`), reproduced from realistic input — a fragment containing a whitespace-only line shorter than the block's indentation; (2) a **206 MB `.git`** caused by auto-committing three ~4.3 MB binaries to the tree on every `master` push via a deploy-key workflow that pushes directly to a protected branch; (3) a **test suite that cannot run in parallel** because it shares a mutable `../test/docs` directory and a process-global resolver cache, forcing a `-p 1` crutch. **Top 3 opportunities:** harden the core with a one-line guard plus fuzz tests (cheap, high payoff); move binary distribution to GitHub Releases and purge the binary history; refactor tests to per-test temp dirs and an injectable cache to unlock parallelism and lift orchestration coverage. There are **no security findings of note** and no hardcoded secrets. A newly-scoped workstream — tightening the heuristic comment filter for the nine target languages, beginning with a concrete Kotlin nested-block-comment defect — is the main correctness investment beyond the quick fixes.

---

## 3. Repo map

**Purpose.** A CLI that (a) **embeds** named code fragments / pattern-matched line ranges from source files into `<embed-code …>` tags in Markdown, or (b) **checks** that already-embedded snippets are up-to-date. Built for TeamDev/SpineEventEngine's Hugo documentation sites (`README.md:6-8`).

**Stack & maturity.** Go 1.22.1 (`go.mod:21`), version 1.2.2 (`main.go:33`). 38 Go files, ~7,856 lines (~5,400 non-test). Dependencies: `doublestar` (globbing), `gobwas/glob` (patterns), `yaml.v3`, Ginkgo/Gomega (tests). 345 commits over two years, primarily one author. **Maturity: established internal product**, not a prototype — and now headed for external distribution.

**Control flow.**
```
main.go ─ orchestration, error aggregation, exit codes
  └─ cli/ ─────────── arg + YAML parsing → Config; validation; build []Configuration
       └─ embedding/ ─ EmbedAll / CheckUpToDate; one Processor per doc file
            └─ parsing/ ─ State-machine parser (Start→…→Finish) over doc lines
                 │         Instruction (the <embed-code> tag) extracts content via:
                 ├─ fragmentation/ ─ split source into #docfragment regions; LRU cache
                 ├─ commentfilter/ ─ strip/keep comments per language & mode
                 └─ indent/ ──────── compute & cut common indentation
  shared: configuration/ (settings), files/ (fs helpers), type/ (YAML list types), logging/
```

| Directory | One-line description |
|---|---|
| `main.go` | Entry point; mode dispatch; multi-config aggregation; error formatting |
| `cli/` | `Config` struct (doubles as CLI args + YAML schema), flag/file parsing, validation |
| `configuration/` | `Configuration` value object + defaults |
| `embedding/` | `Processor` (per-doc) + `EmbedAll`/`CheckUpToDate` orchestration |
| `embedding/parsing/` | State-pattern parser; `Instruction`; `Pattern` (glob); XML tag parsing |
| `embedding/commentfilter/` | Per-language comment stripping (Java/Go/Python/VB/XML/…) |
| `fragmentation/` | `#docfragment` extraction, partitions, generic LRU `cache` |
| `indent/` | Common-indentation detection and trimming |
| `files/`, `type/`, `logging/` | FS existence checks; YAML list unmarshalers; slog handler + `file://` refs |
| `test/resources/` | Java/Kotlin code + Markdown + YAML fixtures |
| `bin/` | **Committed** linux/macos/windows binaries (~13 MB) — *to be removed* |

**Surprises (all [Fact]):** (1) the three platform binaries are committed *and re-committed by CI on every master push* (`.github/workflows/build_binaries.yml:57-63`), inflating `.git` to 206 MB; (2) `logging/logger.go:31` imports the legacy `golang.org/x/net/context` shim instead of stdlib `context`; (3) `fragmentation/partition.go:77-88` uses `panic`/`recover` for an array-bounds check a single `if` would handle.

---

## 4. Audit report

### 4.0 Strengths (preserve these)
- **[Fact] Enforced quality gate.** `.golangci.yml` enables ~40 linters including `gosec`, `gocyclo`, `gocognit`, `cyclop`, `funlen`, `errorlint`, `varnamelen`; CI runs it on every PR across three OSes (`check.yml:24-33`).
- **[Fact] Behavior-asserting tests, 75.9% integrated coverage.** `commentfilter` 91.8%, `embedding` 86.5%, `type` 93.1%, `indent` 100% (own-package). Tests check exact output lines, not just non-error.
- **[Judgment] Clean, extensible parser.** The `State` interface (`state.go:26-43`) + declarative `Transitions` map (`constants.go:34-53`) is readable and open for extension. No circular dependencies (Go enforces this).
- **[Fact] Good diagnostics & typed errors.** `PatternNotFoundError`, `MissingCodeFenceError`, etc. matched via `errors.As` (`processor.go:301-323`) and rendered as clickable `file://…:line` references (`logging/logger.go:99-117`).
- **[Fact] Clean security posture.** No secrets, no `TODO/FIXME`, no XXE (Go's `encoding/xml` resolves no external entities), no injection sinks.

### 4.1 Architecture & design
- **M — [Judgment] `processor.go` mixes two responsibilities.** It holds both the per-file `Processor` (`processor.go:44-186`) and the file-set orchestration `EmbedAll`/`CheckUpToDate`/`requiredDocs`/glob helpers (`processor.go:188-463`). Splitting discovery/orchestration from single-file processing improves cohesion.
- **L — [Judgment] `cli.Config` is overloaded.** One struct serves as CLI-arg holder, YAML schema (`cli.go:64-75`), and carries `Mode`/`ConfigPath`. Pragmatic but couples three concerns.
- **L — [Fact] Pointless wrapper.** `EmbedCodeSamplesResult` (`cli.go:90-92`) embeds `embedding.EmbedAllResult` with no added behavior; `CheckCodeSamples`/`EmbedCodeSamples` (`cli.go:102-118`) are thin pass-throughs.

### 4.2 Code quality
- **H — [Fact, reproduced] Latent panic in `indent.CutIndent`.** `linesChanged[i] = line[redundantSpaces:]` (`indent/indent.go:65`) lacks a guard that `len(line) >= redundantSpaces`. `MaxCommonIndentation` (`indent.go:34`) ignores whitespace-only lines when computing the minimum, so a fragment with an 8-space code line and a 2-space blank line yields `redundantSpaces=8` applied to the length-2 line. Reproduced verbatim: `panic: runtime error: slice bounds out of range [8:2]`. Reachable via `fragment.text` (`fragment.go:72`) and `instruction.matchingLines` (`instruction.go:252,263,289`). Trailing-whitespace blank lines inside indented blocks are common, so this crashes the tool on realistic input (recovered only to a stack-trace exit via `HandlePanic`).
- **M — [Fact] Double file read + dead branch in encoding check.** `IsEncodedAsText` (`encoding.go:31-42`) reads the whole file via `os.ReadFile`, then `DoFragmentation` reads it again with a scanner (`fragmentation.go:103-117`). `areASCIIEncoded` (`encoding.go:47-55`) is fully redundant: ASCII ⊂ UTF-8, so `utf8.Valid` already covers it.
- **M — [Fact] Duplicated joined-error formatting.** `main.go:138-170` (`formatError`/`flattenedErrors`) and `logging/logger.go:173-195` (`formatPanicMessage`) independently implement the same "unwrap `errors.Join` into a bullet list" logic.
- **M — [Fact] Write-only / dead fields.** `Fragmentation.Configuration` (set `fragmentation.go:90`, never read) and `Fragmentation.SourcesRoot` (set `:78`, never read) are dead; `_, err := filepath.Abs(codeRoot.Path)` (`:79`) computes a value only to discard it. `resolver.go:184` passes an empty `config.Configuration{}` precisely because the field is unused. The `unused` linter misses these because the fields are *written*.
- **L — [Fact] `panic`/`recover` for bounds checking.** `safeAccess` (`partition.go:77-88`) uses `recover()` where `index >= 0 && index < len(slice)` suffices.
- **L — [Fact] Swallowed `Close()` error.** `DoFragmentation`'s deferred `err = file.Close()` (`fragmentation.go:108-110`) assigns to a non-named return after the value is computed; the close error is dropped.
- **L — [Fact] Dead defensive branch.** `Pattern.Match` recompiles the glob when `p.matcher == nil` (`pattern.go:111-113`), but `NewPattern` always sets `matcher`.
- **L — [Judgment] Missing defensive length guard.** `IsContentChanged` (`context.go:136-144`) indexes `c.Result[i]` up to `c.lineIndex` (= `len(source)`), relying on an early content mismatch to avoid an out-of-bounds read when an embedding shrinks the file. Holds in practice but is fragile.

### 4.3 Comment-filter reliability *(new workstream — see Decisions)*
[Judgment] Root cause: `commentfilter` is a single-pass byte/marker scanner (`marker_comment_filter.go`) with **no per-language lexer and no cross-line string state** — only *block-comment* state persists across lines (`marker_comment_filter.go:65-77`). This yields concrete, language-specific defects, listed in the owner's priority order:

| # | Language | Current handling | Concrete risk (`file:line`) | Sev |
|---|---|---|---|---|
| 1 | **Kotlin** | C-style markers, all modes | **[Fact] Nested block comments `/* /* */ */` unsupported** — `consumeActiveBlock` ends at the *first* `*/` (`marker_comment_filter.go:123-137`); no depth counter. String templates `"${…}"` / raw `"""…"""` treated as plain quotes. | High |
| 2 | Java | C-style, all | Text blocks `"""…"""` (Java 15+) not modeled. | Med |
| 3 | TypeScript | JS, all | Template literals with `${…}` interpolation and regex literals confuse the scanner; backtick handled as a plain quote (`config.go:25`). | Med-High |
| 4 | JavaScript | JS, all | Same template-literal / regex caveats. | Med-High |
| 5 | C# | C#, all | Verbatim `@"…"` (doubled-quote escaping) and interpolated `$"…"` mis-scanned — scanner assumes `\` escaping (`marker_comment_filter.go:162`). | Med-High |
| 6 | C++ | C-style, regular | Raw string literals `R"(…)"` not handled. | Med |
| 7 | Python | hash, none/all | **Triple-quoted strings spanning lines not tracked** (no cross-line string state) → a `#` on a docstring continuation line is wrongly stripped. | Med-High |
| 8 | Go | Go, regular | Raw strings (backticks, no `\` escaping) slightly mismodeled. | Low |
| 9 | Visual Basic | VB filter, doc modes | `""`-doubled string escaping not handled (`visual_basic.go:57`). | Low |

*All nine languages are already registered (`config.go:29-82`); the issue is robustness, not coverage.* Fix strategy: targeted patches (nested-block depth counter for #1/#2; cross-line string state for #7; verbatim/interpolated handling for #5; template-literal handling for #3/#4) rather than full lexers.

### 4.4 Security
- **[Judgment] Healthy.** One informational item only: an embedding's `file="…"` attribute is `filepath.Join`-ed onto the code root (`resolver.go:173`) with no containment check, so `file="../../secret"` would escape the root. Threat model is local, trusted, self-authored docs → **informational/Low**, not a real boundary.

### 4.5 Testing
- **M — [Fact] Orchestration is 0% covered.** All of `main.go` — `checkByConfigs`, `embedByConfigs`, `formatError`, `flattenedErrors`, `printFiles` (`main.go:82-242`) — is untested; the multi-config aggregation and error-flattening logic has no tests.
- **M — [Fact] Non-trivial logic untested.** The LRU `evictOldest` (`cache.go:108-123`) is 0% — tests never exceed the 100-entry limit (`resolver.go:33`), so the eviction path is unexercised. Error-formatting paths `codeFileReference`/`unresolvedSourceError` (`resolver.go:213-256`) are also 0%.
- **M — [Fact] Parser's own-package coverage is 42.1%.** State-transition files (`blank_line.go:43`, `start.go`, etc.) are exercised only cross-package by `embedding`'s tests, so the parser is less directly tested than it appears. `configuration` has no test file (0%, but trivial).
- **M — [Fact] Shared mutable state forces `-p 1`.** `embedding_test.go:39` uses a shared relative `../test/docs` dir (created/removed per spec, `:61,65`); the resolver cache is a process global cleared manually (`fragmentation_test.go:52`). CI works around this with `go test -p 1` (`check.yml:31-33`).

### 4.6 Performance
- **[Judgment] A non-issue for the workload.** All in-memory over small files; no DB, network, or concurrency. Only items: the double-read/redundant-ASCII check (4.2) and an O(N²) `slices.Contains(p.requiredDocPaths, …)` per doc (`processor.go:106,142,165`) — negligible at realistic doc counts.

### 4.7 Dependencies
- **L — [Fact] Legacy shim.** `golang.org/x/net/context` (`logging/logger.go:31`) is a pre-Go-1.7 relic; stdlib `context` is a drop-in replacement and removes a direct reason to pull `golang.org/x/net`.
- **L — [Fact] Behind latest, no known CVEs.** `doublestar v4.6.1→4.10.0`, `ginkgo 2.20→2.31`, `gomega 1.34→1.42`, `glob v0.2.3` (2018, unmaintained but stable). `go.sum` clean; `-mod=readonly` enforced.

### 4.8 DevEx & operations
- **H — [Fact] 206 MB `.git` from committed binaries.** `bin/` holds three ~4.3 MB binaries, and `build_binaries.yml:57-63` rebuilds and commits them to `master` on every push (23 "Update binaries" commits). A fresh clone pulls 206 MB for <1 MB of source. The workflow uses a deploy key (`build_binaries.yml:27-37`) to push directly to a protected branch — unusual and fragile.
- **L — [Fact] Hand-maintained version.** `const Version = "1.2.2"` (`main.go:33`) is decoupled from any git tag and must be bumped manually.

### 4.9 Documentation
- **L — [Fact] README arg list incomplete.** The "Arguments" section (`README.md:57-63`) omits the `-doc-excludes`, `-info`, and `-stacktrace` flags that `ReadArgs` defines (`cli.go:128-138`).
- **Required — [Fact] No top-level LICENSE.** License terms live only in per-file headers ("All rights reserved" + BSD-style disclaimer). With external distribution decided, a root **Apache-2.0** `LICENSE` (and consistent headers) is now a required deliverable.

---

## 5. Improvement strategy

**Theme 1 — Harden the core against real-world input.** A confirmed panic and a few missing guards sit in the hot path (`indent`, `pattern`, `context`), and the comment filter mishandles common constructs in the target languages. *Target state:* core string/slice operations never panic on valid files; per-language comment/string handling is correct for the nine priority languages, proven by a golden-file corpus. *Principle:* a documentation tool must fail with a diagnostic, never a stack trace, and must never silently corrupt an embedded snippet. *Trade-off:* the comment-filter work is the larger effort; sequence it by the owner's priority order and stop when payoff tapers.

**Theme 2 — Stop shipping binaries through git.** The single largest hygiene problem is operational. *Target state:* binaries distributed via GitHub Releases on tag; `bin/` removed from the tree; the binary history purged; the deploy-key auto-commit workflow retired. *Principle:* build artifacts don't belong in source history. *Trade-off:* consumers fetch from Releases instead of cloning binaries (a documented change); the history purge is destructive and must be scheduled with branch-protection coordination.

**Theme 3 — Make the test suite isolated and parallel-safe.** `-p 1` is a symptom of shared mutable state. *Target state:* each test uses `t.TempDir()`/`os.MkdirTemp`; the resolver cache is injected or deterministically reset; `-p 1` removed. *Principle:* tests own their state.

**Theme 4 — Close coverage gaps in orchestration and error paths.** *Target state:* `main.go` aggregation, cache eviction, and error formatters are unit-tested; duplicated error-formatting logic is unified. *Principle:* the code that runs only on failure is the code you most need tested.

**Theme 5 — Be a credible external dependency.** *Target state:* root Apache-2.0 `LICENSE`, accurate README, versioned Releases with attached binaries. *Principle:* external consumers judge a tool by its packaging and docs as much as its code.

**Definition of done (measurable signals):**
- Zero High findings; fuzz tests for `indent`/`pattern` run in CI with zero crashes.
- `go test ./...` passes **without `-p 1`**.
- Integrated coverage ≥ 80%, with `main.go` and `cache.evictOldest` > 0%.
- `bin/` absent from the tree; `.git` < 5 MB after purge; binaries attached to the latest GitHub Release.
- Root Apache-2.0 `LICENSE` present; README arg list matches `ReadArgs`.
- Golden-file comment-filter corpus passes for Kotlin → … → Go (Visual Basic best-effort).

**Explicitly NOT doing (calibrated to maturity):** migrating off Ginkgo (works well; churn not worth it); restructuring `cli.Config`'s dual role (pragmatic as-is); chasing 100% coverage (diminishing returns past ~85%); writing full per-language lexers where a targeted fix suffices.

---

## 6. Task plan

### Quick wins (high impact, S effort — do immediately)
- **QW1** Guard `CutIndent` against short lines (clamp/skip when `len(line) < redundantSpaces`) — fixes the confirmed panic. `indent/indent.go:60-70`.
- **QW2** Replace `golang.org/x/net/context` with stdlib `context`. `logging/logger.go:31`.
- **QW3** Delete dead fields/branches: `Fragmentation.Configuration`/`SourcesRoot` + discarded `Abs` (`fragmentation.go:78-90`); collapse `areASCIIEncoded` into `utf8.Valid` (`encoding.go`); drop the dead `matcher==nil` branch (`pattern.go:111`).
- **QW4** Add `-doc-excludes/-info/-stacktrace` to the README arg list (`README.md:57-63`).

### Task table

| ID | Title | Files | Acceptance criteria | Effort | Change risk | Deps |
|---|---|---|---|---|---|---|
| **M0-1** | Fuzz/table tests for `indent` & `pattern` | `indent/*_test.go`, `parsing/pattern_test.go` (new) | Fuzz targets in CI; whitespace/short-line/multiline cases asserted | M | Low | — |
| **M0-2** | Test harness uses temp dirs; remove `-p 1` | `embedding/embedding_test.go`, `check.yml:33` | Suite passes with default parallelism; no shared `../test/docs` | M | Med | — |
| **M0-3** | Injectable/reset resolver cache | `fragmentation/resolver.go:42`, `cache.go` | Cache resettable per test; no package-global reliance | M | Med | M0-2 |
| **M1-1** | Fix `CutIndent` panic (QW1) | `indent/indent.go:65` | Reproduction input returns trimmed lines, no panic; regression test | S | Low | M0-1 |
| **M1-2** | Defensive guard in `IsContentChanged` | `context.go:136-144` | Shrinking-embed case covered; no OOB | S | Low | M0-2 |
| **M2-1** | Binaries → GitHub Releases | `build_binaries.yml`, `bin/`, `.gitignore`, README | Tagged release attaches 3 binaries; `bin/` untracked; deploy-key push removed | M | Med | — |
| **M2-1b** | **Purge `bin/` history** | whole repo history | `git filter-repo --path bin --invert-paths`; `.git` < 5 MB; tags re-pushed | M | **High** | M2-1 |
| **M2-2** | Split orchestration out of `processor.go` | `embedding/processor.go` | Discovery/orchestration in own file; behavior unchanged | M | Med | M0-2 |
| **M2-3** | Unify joined-error formatting | `main.go:138-170`, `logging/logger.go:173-195` | Single shared helper; tests for single + joined errors | S | Low | — |
| **M-LIC** | Add Apache-2.0 LICENSE + headers | `LICENSE` (new), file headers, README | Root `LICENSE` = Apache-2.0; headers consistent; README states license | S | Low | — |
| **M3-1** | Tests for `main.go` aggregation & `printFiles` | `main_test.go` (new) | `checkByConfigs`/`embedByConfigs`/formatters > 0% covered | M | Low | M0-2 |
| **M3-2** | Test cache LRU `evictOldest` | `fragmentation/cache_test.go` (new) | Eviction-at-limit asserted; `evictOldest` covered | S | Low | — |
| **M3-3** | Single-read fragmentation | `encoding.go`, `fragmentation.go` | Source read once; encoding check folded in/removed | M | Med | M0-1 |
| **M3-4** | Replace `panic/recover` bounds check | `partition.go:77-88` | Plain bounds check; tests green | S | Low | — |
| **M-CF-1** | Kotlin: nested block comments + string templates | `commentfilter/*` | Golden-file corpus for nested `/* */` and `"""…"""` passes all modes | M | Med | M0-1 |
| **M-CF-2** | Java text blocks; TS/JS template literals & regex | `commentfilter/*` | Golden files for Java `"""`, TS/JS `${…}` / regex pass | M | Med | M-CF-1 |
| **M-CF-3** | C# verbatim/interpolated; C++ raw strings; Python triple-quotes | `commentfilter/*` | Cross-line string state added; golden files pass | L | Med | M-CF-2 |
| **M-CF-4** | Go raw strings; VB doubled-quote escaping | `commentfilter/*`, `visual_basic.go` | Golden files pass (best-effort) | S | Low | M-CF-3 |

### Milestones
- **Milestone 0 — Safety net:** M0-1, M0-2, M0-3 (tests + parallel-safety before refactors).
- **Milestone 1 — Critical / correctness:** M1-1, M1-2 (plus QW1).
- **Milestone 2 — High-leverage / packaging:** M2-1, M2-1b (gated), M2-2, M2-3, M-LIC.
- **Milestone 3 — Quality & polish:** M3-1…M3-4 + remaining quick wins.
- **Workstream M-CF — Comment-filter reliability:** M-CF-1…M-CF-4, in owner priority order, runs in parallel after Milestone 0.

### Top-3 implementation sketches

**M1-1 — Fix `CutIndent` panic.** *Approach:* clamp the slice to the line length. In `CutIndent` (`indent.go:63-67`), bound the cut by `min(redundantSpaces, len(line))`. *Steps:* add the bound; add a regression test using the reproduced fixture (8-space code line + 2-space blank line); add a fuzz target over random indented blocks. *Gotcha:* leave the normal case (lines ≥ indentation) unchanged; verify `indent_test.go` (100% today) still passes.

**M2-1 / M2-1b — Binaries → Releases, then purge.** *Approach (M2-1):* trigger on `push: tags: ['v*']`; build matrix; upload via `softprops/action-gh-release`; `git rm -r --cached bin/`; add `bin/` to `.gitignore`; drop the `ssh-agent`/deploy-key + `add-and-commit` steps (`build_binaries.yml:27-37,57-63`); point README at Releases. *Approach (M2-1b):* on a scheduled maintenance window, `git filter-repo --path bin --invert-paths`, force-push all refs and tags after temporarily lifting `master` protection. *Gotchas:* M2-1b rewrites every SHA — announce ahead, require all forks/clones to re-clone, and have a maintainer (not this automation) perform the force-push. M2-1 halts growth; only M2-1b reclaims space.

**M-CF-1 — Kotlin nested block comments.** *Approach:* give the C-family block scanner a nesting depth counter so `/* /* */ */` only closes at depth 0, and add cross-line *string* state (mirroring the existing block state) so `"""…"""` raw strings and `"${…}"` templates protect their contents. *Steps:* extend `blockState` with a `depth` and a `inString` mode; update `consumeActiveBlock` (`marker_comment_filter.go:118-140`) and `consumeQuotedSegment` (`:143-152`); build a golden-file corpus under `test/resources` covering nested comments and multiline strings for every mode. *Gotcha:* enable nesting only for languages that actually allow it (Kotlin) — Java/C/C++ do **not** nest block comments, so keep nesting per-syntax to avoid regressions.

---

## 7. Appendix — evidence

**Build & tests:** `go build ./...` and `go vet ./...` clean. Full suite passes (`go test ./... -p 1`).

**Integrated coverage = 75.9%** (`go test ./... -coverpkg=./...`). Per-package, own-tests:

| Package | Coverage | | Package | Coverage |
|---|---|---|---|---|
| `indent` | 100.0% | | `cli` | 75.8% |
| `commentfilter` | 91.8% | | `fragmentation` | 68.8% |
| `type` | 93.1% | | `parsing` | 42.1% * |
| `embedding` | 86.5% | | `logging` | 21.8% |
| `files` | 78.6% | | `configuration` | 0.0% (trivial) |
| | | | `main` | 0.0% |

\* The parser is exercised mainly cross-package by `embedding`'s tests, so its true coverage is higher than the own-package figure.

**Reproduced panic (`indent.CutIndent`)** — verbatim copies of `MaxCommonIndentation` and `CutIndent` fed a fragment `["        System.out.println(\"Hi\");", "  ", "        return;"]`:
```
computed common indentation = 8
calling CutIndent...
panic: runtime error: slice bounds out of range [8:2]
```

**Repository bloat:** `.git` = 206 MB; 23 "Update binaries" commits; `bin/` = three binaries ~4.3 MB each.

---

*Prepared with Claude Code. Analysis-only audit — no source code was modified in its production.*
