# Coding Style

## Coding Style & Naming Conventions
- Use object-oriented design with structs and methods; avoid large functions and global state.
- Follow generally accepted object oriented design patterns.
- Keep Go code `gofmt`-clean; package names should be short and lower-case.
- Prefer focused edits over broad refactors unless required.
- Add and maintain explanatory comments for non-obvious logic.
- Use `git mv` for file moves/renames to preserve history.

## Diff Discipline (Required)
- Keep changes minimal and directly tied to the user's request or locked decision.
- Do not rewrite files from scratch, arbitrarily rewrap lines, reorder unrelated sections, or normalize prose style unless the user explicitly asks for that cleanup.
- This applies to `AGENTS.md`, slides, white papers, drafts, TODO files, code, and all other repo files.
- For public prose artifacts, edits must be traceable to a concrete user request or locked decision, not general polish.

## Runtime Artifact Hygiene (Required)
- Never put temporary test files, Go cache directories, build caches, or other runtime artifacts in this repo.
- Runtime artifacts must go under `/tmp` subdirectories.
- Configure tools with temp/cache paths under `/tmp` when they would otherwise write into the working tree.

## Error Handling Policy (Required)
- Never use `|| true` in scripts, templates, or make recipes. Always inspect
  command exit codes explicitly with `if/else` branches and handle each outcome.
- For non-fatal cleanup/diagnostics steps, record command status (exit code and
  logs) explicitly; do not fail silently.
- In Go code, never ignore errors with `_ = ...`; handle, propagate, or report
  errors explicitly.
- Run `errcheck ./...` and keep it passing for Go changes.
