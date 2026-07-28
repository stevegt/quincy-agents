# TODO-jusuk - Mogent: Modular Agent Prompt Manager

A CLI tool for assembling AGENTS.md files from modular, scoped templates.

## Decision Intent Log

ID: DI-jusuk
Date: 2026-07-24 12:00:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Build mogent as a Go CLI with tag-based scoping, slash hierarchy, header attribute tags, grouped TOML config, and section-aware semantic diff. Design for promise grid integration but implement later.
Intent: Create a lightweight but flexible tool for managing agent prompts across repos, orgs, and people. Enable cross-org sharing via promise grid in future.
Constraints: Must work with existing AGENTS.md format. Must be simple to adopt incrementally.
Affects: tools/mogent/, AGENTS.toml, .mogent/

ID: DI-jusuk-edit
Date: 2026-07-24 13:00:00
Status: active
Author: 95124070+Qu1ncyRy4n@users.noreply.github.com (Quincy Ryan)
Decision: Edit command opens module in $EDITOR or assembled AGENTS.md. No special diff/merge on save.
Intent: Keep edit simple - use existing editor. Future versions can add merge assistance.
Constraints: Requires $EDITOR to be set.
Affects: tools/mogent/cmd/edit.go

## Subtasks

- [x] jusuk.1 Project scaffolding - Go module, CLI skeleton, basic build
- [x] jusuk.2 TOML parser - Parse AGENTS.toml with categories, tags, sources
- [x] jusuk.3 Module parser - Parse markdown with header attribute tags {#tag1 #tag2}
- [x] jusuk.4 Tag resolver - Hierarchical tag matching with parent-implies-children
- [x] jusuk.5 Assembly engine - Combine modules based on active scopes
- [x] jusuk.6 Build command - mogent build outputs AGENTS.md
- [x] jusuk.7 Edit command - mogent edit <module> and mogent edit (assembled)
- [x] jusuk.8 Diff command - Section-aware diff between scopes
- [x] jusuk.9 List command - List available and active modules
- [x] jusuk.10 Dogfood - Use mogent on this repo's AGENTS.md

## Design References

- docs/examples/DESIGN-SUMMARY.md - Final design decisions
- docs/examples/AGENTS-style-e-hierarchical.md - TOML style reference
- docs/examples/AGENTS-style-f-slashes.md - Slash notation reference
- docs/thought-experiments/TE-tavim-mogent-module-reference-model.md - Follow-up TE for block references, presets, metadata, render order, and diff model. Manual handle allocation approved by user because `tools/mint-handle` was unavailable at the approved paths.
