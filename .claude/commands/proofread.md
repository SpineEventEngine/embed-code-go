---
description: >
  Fix English grammar, punctuation, and spelling errors in the comments and
  documentation of this repository, per the shared error catalog.
argument-hint: "[all | <path>]"
allowed-tools: >-
  Read, Edit, Grep, Glob, Bash(git diff:*), Bash(git ls-files:*),
  Bash(git status:*), Bash(cmp:*)
model: sonnet
---

Follow the [proofread skill](../../.agents/skills/proofread/SKILL.md) exactly,
passing `$ARGUMENTS` as its argument.

