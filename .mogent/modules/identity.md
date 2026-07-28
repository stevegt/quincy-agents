# Identity

<!--
agent_module:
  id: identity
  tldr: Defines the assistant role and project context.
-->
You are a software engineering assistant helping with a Go project. Your role is to help with coding, debugging, and architectural decisions while following the project's conventions.

## Project Overview

<!--
agent_module:
  id: project-overview
  tldr: Describes the project and its decision-first workflow.
-->
This is a Go project using modular agent prompt management. The project follows decision-first development where decisions are locked before implementation.

## Tech Stack

<!--
agent_module:
  id: tech-stack
  tldr: Lists the main language, build, lint, and version-control tools.
-->
- Language: Go
- Build: go build, go test
- Linting: gofmt, errcheck
- Version control: Git

## Project Structure

<!--
agent_module:
  id: project-structure
  tldr: Captures repository layout and source-of-truth rules.
-->
- Keep prose drafts, presentation decks, and small repo tools in purpose-named top-level paths; avoid `internal/` and `pkg/`.
- Keep planning artifacts in the root `TODO/` directory. Optional root `DR/` and `docs/thought-experiments/` directories hold DR and TE records when needed.
- Do not commit local state files (for example `.grok`, `.grok.lock`) or generated binaries.
- Run Go tests from the relevant tool module, for example `tools/mint-handle/`.
- Keep Go code `gofmt`-clean.
- `AGENTS.md` is the canonical home for repo-wide workflow, review, and vocabulary rules.
- Keep repo-wide rules here; keep private runtime notes out of committed files.
