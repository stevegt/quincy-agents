# Constraints

<!--
agent_module:
  id: constraints
  tldr: Defines hard prohibitions and always-on requirements.
-->
## Never Do

<!--
agent_module:
  id: never-do
  tldr: Lists forbidden actions and patterns.
-->
- Do not commit secrets, tokens, or credentials
- Do not commit generated binaries or local state files
- Do not use `|| true` in scripts or make recipes
- Do not ignore errors with `_ = ...` in Go code
- Do not put runtime artifacts in the repo (use /tmp)
- Do not force-push unless explicitly authorized
- Do not open PRs unless explicitly asked
- Do not remove comments or documentation; update them if outdated or incorrect
- Do not assume defaults for locked categories unless the user explicitly approves

## Always Do

<!--
agent_module:
  id: always-do
  tldr: Lists mandatory engineering and compliance practices.
-->
- Handle errors explicitly in Go code
- Run `errcheck ./...` and keep it passing
- Keep Go code gofmt-clean
- Use %w for error wrapping
- Inspect command exit codes explicitly
- Ask decision questions up front in a single round
- Lock decisions before coding
- Record decisions as Decision Intent entries in TODO files
- Treat user decisions as authoritative and implement to those decisions
- Run a compliance self-review before finalizing
- Fix all non-compliance before handoff
