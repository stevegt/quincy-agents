# Format

## Coding Style

- Use object-oriented design with structs and methods; avoid large functions and global state
- Follow generally accepted object oriented design patterns
- Keep Go code `gofmt`-clean; package names should be short and lower-case
- Prefer focused edits over broad refactors unless required
- Add and maintain explanatory comments for non-obvious logic

## Diff Discipline

- Keep changes minimal and directly tied to the user's request or locked decision
- Do not rewrite files from scratch, arbitrarily rewrap lines, reorder unrelated sections, or normalize prose style unless the user explicitly asks for that cleanup
- This applies to `AGENTS.md`, slides, white papers, drafts, TODO files, code, and all other repo files
- For public prose artifacts, edits must be traceable to a concrete user request or locked decision, not general polish

## Error Handling

- Never use `|| true` in scripts, templates, or make recipes
- Always inspect command exit codes explicitly with `if/else` branches and handle each outcome
- For non-fatal cleanup/diagnostics steps, record command status (exit code and logs) explicitly; do not fail silently
- In Go code, never ignore errors with `_ = ...`; handle, propagate, or report errors explicitly

## Runtime Artifact Hygiene

- Never put temporary test files, Go cache directories, build caches, or other runtime artifacts in this repo
- Runtime artifacts must go under `/tmp` subdirectories
- Configure tools with temp/cache paths under `/tmp` when they would otherwise write into the working tree

## Decision Lock

- Produce a Decision Lock summary with decision IDs before code edits begin
- Do not proceed if any required decision is missing, ambiguous, or conflicting
- Stop and ask immediately if a new decision need appears during implementation

## Required Handoff Artifacts

- `Decision Compliance: PASS/FAIL`
- Decision Matrix mapping each locked decision ID to implementation evidence
- Inline diff annotations in the form `path:line -> decision_id -> rationale`
- Runtime Path Touch Matrix listing each approved runtime path/pattern, action used, and where it is implemented/validated
- `Exceptions:` listing only user-approved deviations
- Every non-trivial behavior change must include intent provenance per existing DI requirements

## Glossary

- **TE**: Thought Experiment. Analysis doc under `docs/thought-experiments/TE-<handle>-<slug>.md`. The handle is a proquint minted by `tools/mint-handle`.
- **DR**: Decision Request. Open question or decision-tracking record under `DR/DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`.
- **DI**: Decision Intent. Locked decision record inside a `## Decision Intent Log` in a TODO file. New DI ID format is `DI-<handle>`, where `<handle>` is minted by `tools/mint-handle`.
- **DF**: Decision Framing. The multiple-choice intake round used to lock a decision after any required TE narrows the alternatives.
- **TODO**: Task-tracking file under `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. Existing numeric TODO files are legacy records until migrated.
- **pCID**: Protocol CID. A pCID is the content hash of a spec document that defines a wire protocol; it is analogous to a TCP/UDP port number with no central registry because the spec hash is the port number. A pCID is not the hash of a particular message, payload, or promise body.
