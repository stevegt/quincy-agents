# Instructions

<!--
agent_module:
  id: instructions
  tldr: Defines the default workflow and repository process rules.
-->
## Workflow

<!--
agent_module:
  id: workflow
  tldr: Gives the basic work loop for coding tasks.
-->
1. Understand the request before coding
2. Ask clarifying questions if the request is ambiguous
3. Make minimal, focused changes
4. Test your changes before handoff
5. Follow existing patterns in the codebase

## Code Changes

<!--
agent_module:
  id: code-changes
  tldr: Sets expectations for focused edits and comment preservation.
-->
- Make minimal, focused changes directly tied to the request
- Do not rewrite files from scratch unless explicitly asked
- Preserve existing code comments unless replacing with better ones
- Use `git mv` for file moves/renames to preserve history

## Testing

<!--
agent_module:
  id: testing
  tldr: Defines deterministic Go testing and runtime artifact placement.
-->
- Use Go's standard `testing` package with deterministic tests
- Avoid network calls in tests unless explicitly required and documented
- Put test temp files and Go caches under `/tmp`, not inside this repo

## Commits

<!--
agent_module:
  id: commits
  tldr: Defines commit behavior and forbidden commit contents.
-->
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message
- Use short, imperative, capitalized commit subjects
- Summarize changes per file in commit bodies
- Stage files explicitly (avoid `git add .` / `git add -A`)
- Do not open GitHub pull requests unless the user explicitly asks for one
- Do not force-push repo branches unless the user explicitly authorizes the exception
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets

## Decision Protocol

<!--
agent_module:
  id: decision-protocol
  tldr: Requires decisions to be identified, asked, locked, and recorded.
-->
When making non-trivial changes:
1. Identify what decisions need to be made
2. Ask decision questions up front in a single round
3. Use multiple-choice questions when practical
4. Lock decisions before coding
5. Record decisions as Decision Intent entries in TODO files

## Thought Experiments

<!--
agent_module:
  id: thought-experiments
  tldr: Defines pre-decision scenario analysis for non-trivial decisions.
-->
Before locking any non-trivial decision:
1. Run a thought experiment if multiple plausible designs remain
2. Evaluate the decision across multiple concrete scenarios
3. Include: normal operation, failure/corruption, concurrent actors, long-horizon evolution, trust-boundary changes, scale effects
4. Compare alternatives under the same assumptions
5. State what each alternative makes easier, harder, and what new obligations it creates

## DR/DI Protocol

<!--
agent_module:
  id: dr-di-protocol
  tldr: Defines decision records as the source of truth.
-->
- DR and DI logs are the primary source of truth for decisions and open questions
- Documents and code are outputs of that process and must link back to DR/DI records
- Person identity in DR/DI records must use full email with label format: `user@example.com (FirstName)`
- A settled statement in docs (or critical logic in code comments) must cite at least one DI ID
- An unresolved question or uncertainty must cite at least one DR ID
- If an unresolved question has no DR yet, create a DR before finalizing the change

## Comment Preservation

<!--
agent_module:
  id: comment-preservation
  tldr: Preserves explanatory code comments and requires DI provenance.
-->
- Never remove existing code comments unless they are replaced in the same patch by equal-or-better explanatory comments near the same logic
- When rewriting or refactoring code, port old explanatory intent first, then improve wording
- If a touched non-trivial code block has no comments, add explanatory comments
- Do not treat shorter comments as better unless they preserve all important intent
- For any non-trivial behavior change, include a behavior-level comment with:
  - `Intent:` a short, clear rationale
  - `Source:` a DI ID in the format `DI-<handle>`

## Public Artifacts

<!--
agent_module:
  id: public-artifacts
  tldr: Keeps public artifacts clean while preserving provenance in references.
-->
- Do not put DI, DR, TODO, or TE references in slides; keep public decks visually clean
- In white papers, use GitHub Markdown footnote tags for DI, DR, TODO, or TE references
- White papers with DI, DR, TODO, or TE footnotes must include a `## References` section at the bottom of the document
- The `## References` section must cite the appropriate repo files using relative file paths

## TODO Tracking

<!--
agent_module:
  id: todo-tracking
  tldr: Defines TODO index and task-file conventions.
-->
- Maintain a `./TODO/` directory for tracking tasks and plans
- Maintain a `./TODO/TODO.md` file that lists small tasks and the other TODO files
- New TODO files use `TODO-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`
- Existing numeric TODO files remain valid legacy records until migrated
- Sort `TODO.md` by priority, not filename
- When completing a TODO, mark it as done by checking it off
- Within a TODO file, include numbered checkboxes for subtasks where helpful
