# Testing & Commits

## Testing Guidelines
- Use Go's standard `testing` package with deterministic tests.
- Avoid network calls in tests unless explicitly required and documented.
- Put test temp files and Go caches under `/tmp`, not inside this repo.

## Commit & Pull Request Guidelines
- Treat a line containing only `commit` as: add and commit all changes with an AGENTS-compliant message.
- Use short, imperative, capitalized commit subjects.
- Summarize changes per file in commit bodies.
- Stage files explicitly (avoid `git add .` / `git add -A`).
- Do not open GitHub pull requests unless the user explicitly asks for one.
- Do not force-push repo branches unless the user explicitly authorizes the exception.
- Do not commit local state files, generated binaries, credentials, tokens, signing keys, or other secrets.

## Glossary
- **TE**: Thought Experiment. Analysis doc under `docs/thought-experiments/TE-<handle>-<slug>.md`. The handle is a proquint minted by `tools/mint-handle`.
- **DR**: Decision Request. Open question or decision-tracking record under `DR/DR-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`.
- **DI**: Decision Intent. Locked decision record inside a `## Decision Intent Log` in a TODO file. New DI ID format is `DI-<handle>`, where `<handle>` is minted by `tools/mint-handle`.
- **DF**: Decision Framing. The multiple-choice intake round used to lock a decision after any required TE narrows the alternatives.
- **TODO**: Task-tracking file under `TODO/TODO-<handle>-<slug>.md`, where `<handle>` is minted by `tools/mint-handle`. Existing numeric TODO files are legacy records until migrated.
- **pCID**: Protocol CID. A pCID is the content hash of a spec document that defines a wire protocol; it is analogous to a TCP/UDP port number with no central registry because the spec hash is the port number. A pCID is not the hash of a particular message, payload, or promise body.
