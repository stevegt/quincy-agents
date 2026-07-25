# Project Structure & Build

## Project Structure & Module Organization
- Keep prose drafts, presentation decks, and small repo tools in purpose-named top-level paths; avoid `internal/` and `pkg/`.
- Keep planning artifacts in the root `TODO/` directory. Optional root `DR/` and `docs/thought-experiments/` directories hold DR and TE records when needed.
- Do not commit local state files (for example `.grok`, `.grok.lock`) or generated binaries.

## Build, Test, and Development Commands
- Run Go tests from the relevant tool module, for example `tools/mint-handle/`.
- Keep Go code `gofmt`-clean.

## Agent Instruction Architecture (Required)
- `AGENTS.md` is the canonical home for repo-wide workflow, review, and vocabulary rules.
- Keep repo-wide rules here; keep private runtime notes out of committed files.
