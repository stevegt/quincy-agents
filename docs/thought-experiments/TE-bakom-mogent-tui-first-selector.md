# TE-bakom - Mogent TUI-first selector

TE ID: TE-bakom

## Decision Under Test

Whether the next mogent iteration should keep extending prompt-driven `init` and tree `list`, or introduce a Bubble Tea-backed TUI selector before expanding the agent module library.

## Assumptions

- The current text list is readable but not enough for navigating many selectable heading blocks.
- Tags, global/local module storage, rich diff, and manual `AGENTS.md` drift detection remain unsettled and should not block a selector prototype.
- The user wants fast proof-of-concept dogfooding over backwards compatibility.
- The current repo has intentionally dirty dogfood state in `AGENTS.md`, `AGENTS.toml`, and `AGENTS.md.bak`; this TE should not overwrite that state by itself.
- `tools/mint-handle` is unavailable, so `TE-bakom` is manually assigned under the existing user-approved exception.

## Alternatives

1. Keep improving plain CLI list/init prompts.
2. Build a minimal custom terminal interface with raw ANSI input handling.
3. Add a Bubble Tea-backed `mogent tui` browser that initially views and toggles module blocks in memory, then later writes config.
4. Build module-library/templates first and defer the TUI.

## Scenario Analysis

### Normal Operation

Alice opens a repo with four modules today and thirty modules later. Plain CLI list output can show state but forces Alice to mentally map IDs into TOML edits. A custom ANSI interface would be fast to start but creates terminal-state obligations immediately. Bubble Tea gives Alice a navigable cursor, width-aware view, and predictable key handling while keeping the first iteration small. Building templates first gives more content but worsens the navigation problem the user already hit.

### Failure Or Incomplete Writes

Bob exits midway through selection. Plain CLI and a view-only Bubble Tea prototype avoid partial config writes. A custom raw terminal implementation risks leaving terminal state damaged unless carefully restored. A write-capable TUI must later use atomic writes and backup/drift policy, but that can be deferred if the first command is inspect/toggle-only.

### Concurrent Actors Or Mixed Versions

Carol edits `AGENTS.toml` while Dave runs the selector. Plain CLI and a read-only TUI do not introduce new concurrency problems. A write-capable TUI would need stale-file detection before saving. That belongs with the later manual drift TE rather than the first shell.

### Long-Horizon Evolution

Ellen later wants presets, profile templates, xor groups, local/global modules, and import from manually edited `AGENTS.md`. Plain prompts do not scale to that workflow. A Bubble Tea model can grow panels for categories, modules, blocks, preview, and diagnostics while reusing the same indexing and validation code. A hand-rolled ANSI tool would accumulate UI debt quickly.

### Trust Boundary Changes

Frank eventually points mogent at shared module repositories. A TUI that surfaces source path, module metadata, and selected blocks can make trust boundaries visible. Template-first work without a UI risks hiding provenance in generated files.

### Scale Effects

At larger scale, the important cost is operator attention, not CPU. Bubble Tea adds dependency weight but reduces selector complexity. The repo already uses Go and Cobra; adding a standard Go TUI dependency is a reasonable POC cost if the first integration is isolated.

## Conclusions

Rejected:

- Plain CLI-only iteration, because it does not address the user's current blocker.
- Custom ANSI TUI, because terminal behavior would become a distraction.
- Module-library-first work, because content growth makes the navigation problem worse.

Surviving alternative:

- Add a Bubble Tea-backed `mogent tui` command as a first selector shell. Initial scope: load `AGENTS.toml`, parse configured modules, display a yazi-like tree/list with cursor navigation, show selected/inherited/inactive markers, and support in-memory toggling for selected blocks. Config writing, presets, tags, drift handling, and global/local module policy remain future decisions.

## Implications For Open TODOs And Pending DIs

- `jusuk.9` and `jusuk.10` are extended by a TUI dogfood path.
- A new DI should lock the TUI shell scope and dependency choice before implementation.
- Future TEs remain needed for manual `AGENTS.md` drift handling and local-vs-global module storage.
