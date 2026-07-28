# TE: Mogent Module Reference Model

TE ID: TE-tavim

## Decision Under Test

How should mogent identify, select, assemble, and diff reusable AGENTS.md prompt modules as the module database grows beyond four coarse files?

This TE tests whether mogent should standardize on file-level modules, heading-level blocks, Obsidian-style links, generated block IDs, content hashes, metadata comments, TOML ordering, presets, and a revised diff model.

Related TODO: `jusuk` / follow-up design after `jusuk.8` and `jusuk.10`.

## Assumptions

- Users want to browse a module database in a yazi-like selector.
- Users want recurring complex selections saved as presets.
- Users want TOML to remain the source of render ordering.
- Tags are deferred for now and should not be required for the first version of this redesign.
- The final rendered `AGENTS.md` should not include builder-only metadata.
- The same module database should support small repos, learning projects, research-heavy projects, strict process repos, and fast-iteration repos.
- Manual handle allocation was approved by the user because the required `tools/mint-handle` command was unavailable at the approved paths.

## Alternatives

### Alternative A: File-Level Modules Only

Each selectable unit is a Markdown file.

Example:

```text
modules/
  basic/
    identity-agent-builder.md
  cognition/
    think-hard.md
    think-light.md
```

Presets reference files:

```toml
include = [
  "basic/identity-agent-builder.md",
  "cognition/think-light.md",
]
```

### Alternative B: Heading-Level Blocks in Larger Files

Each selectable unit is a heading block inside a category file. Files are organizational containers.

Example:

```text
modules/
  basic/identity.md
  cognition/thinking-style.md
```

Presets reference blocks:

```text
[[basic/identity#agent-builder]]
[[cognition/thinking-style#think-light]]
```

### Alternative C: Support Both File-Level Modules and Heading-Level Blocks

Users may select either a whole file or a heading block. Presets can mix both.

Example:

```text
[[basic/identity]]
[[cognition/thinking-style#think-light]]
```

### Alternative D: Registry-First Blocks

Every selectable unit is registered in TOML with a stable ID, path, and heading anchor.

Example:

```toml
[[block]]
id = "think-light"
source = "cognition/thinking-style.md"
anchor = "think-light"
```

Presets reference registry IDs:

```toml
include = ["agent-builder", "think-light"]
```

## Metadata Options

### Option 1: Obsidian Comments

```md
%% tldr: Prefer fast practical answers. %%
%% conflicts: cognition/thinking-depth %%
## Think Light
```

This is familiar to Obsidian users and visually compact, but it is less universal outside Obsidian.

### Option 2: HTML Comments with YAML-Like Metadata

```md
<!--
mogent:
  id: think-light
  tldr: Prefer fast practical answers.
  group: cognition/thinking-depth
-->
## Think Light
```

This is valid Markdown, ignored by normal renderers, compatible with Obsidian, and easier to recognize as tool-owned metadata.

### Option 3: TOML-Only Metadata

All metadata lives in `AGENTS.toml`; Markdown contains only final prompt text.

This keeps prose clean but makes modules harder to browse and move because the description, conflicts, and ID are separated from the text they describe.

## Scenario Analysis

### Scenario 1: Normal Operation

Alice builds an AGENTS.md for this Go repo. She wants fast iteration now, stricter design later, Go style, decision process rules, and a lightweight communication style.

- Alternative A makes browsing simple when each option is a small file. It creates many tiny files once dev involvement, thinking style, docs style, and language idioms are decomposed into options.
- Alternative B keeps related options together. Alice can open `cognition/thinking-style.md` and compare `Think Hard`, `Think Light`, and `Socratic` in one place. The selector can still flatten headings into selectable items.
- Alternative C gives maximum flexibility but forces the builder, diff, and selector to explain two selection semantics.
- Alternative D gives strong machine stability but makes simple authoring feel registry-heavy.

Consequence: Alternative B best matches the desired selector behavior if the project accepts one standard selectable unit: the block.

### Scenario 2: Heading Rename and Prose Editing

Bob edits `## Think Light` to `## Fast Practical Reasoning`.

- Alternative A is unaffected if the filename is stable.
- Alternative B breaks Obsidian-style links if the link target is derived only from heading text.
- Alternative C inherits the same breakage for block references.
- Alternative D survives because the stable ID is separate from the heading.

Consequence: If block references use Obsidian-looking syntax, the anchor should resolve to explicit `mogent.id`, not raw heading text. Human-facing headings may change; selected IDs should not.

### Scenario 3: Content Hash IDs

Carol proposes SHA-256 content hashes for block IDs.

- Content hashes guarantee that a reference names exact content.
- They change whenever wording changes, including typo fixes.
- They make presets noisy and hard to review.
- They are useful as integrity metadata but poor as primary human-authored IDs.

Consequence: content hashes should not be primary block IDs. They may be useful later as optional lockfile data for reproducible builds.

### Scenario 4: Failure, Corruption, and Incomplete Writes

Dave has a preset that references `[[cognition/thinking-style#think-light]]`, but the metadata block is missing or duplicated.

- Alternative A can fail on a missing file.
- Alternative B needs validation for missing files, duplicate IDs, missing IDs, and ambiguous anchors.
- Alternative C needs all Alternative B validation plus rules for whole-file selection.
- Alternative D can validate the registry before reading content, but still needs source and anchor validation.

Consequence: heading-level blocks create validation obligations. The builder should fail loudly on missing or duplicate IDs rather than silently dropping content.

### Scenario 5: Concurrent Actors and Mixed Versions

Alice uses a new builder that understands block metadata. Bob uses an older builder that only includes files by TOML module source.

- Alternative A is most compatible with the current implementation.
- Alternative B requires a migration path. Old builders may render whole files and include mutually exclusive options together.
- Alternative C can preserve old file-level behavior while new presets reference blocks, but that compatibility costs complexity.
- Alternative D requires all users to adopt the registry format before presets are useful.

Consequence: the migration should be explicit. A compatibility mode may be acceptable temporarily, but the durable standard should avoid long-term dual semantics.

### Scenario 6: Long-Horizon Evolution

Ellen builds a large personal module database for Godot/Rust learning, Obsidian note maintenance, Go tools, promise-grid apps, neuro analysis, and teacherbot work.

- Alternative A scales by filesystem search but scatters comparisons across many files.
- Alternative B scales by category files and keeps choice families together.
- Alternative C scales poorly in documentation because every command must define when a file reference means "all blocks" versus "this document as one module."
- Alternative D scales for tooling but makes casual manual edits more formal.

Consequence: B with explicit metadata has the best authoring and browsing balance. D's registry strengths can be recovered later by generating an index from block metadata.

### Scenario 7: Trust Boundary Changes

Frank imports a shared module pack from another person.

- Alternative A needs file-path validation and path traversal protections.
- Alternative B also needs block metadata validation and should treat metadata as untrusted input.
- Alternative C has the largest attack and ambiguity surface.
- Alternative D can centralize validation but still depends on trusted parser behavior.

Consequence: all alternatives need safe path resolution. The current code should be changed before this feature ships because it silently skips unreadable modules instead of surfacing errors.

### Scenario 8: Scale Effects

A module database grows to hundreds of blocks.

- Alternative A creates many files and may be harder to rename or reorganize cleanly.
- Alternative B keeps fewer files but requires a block index for fast selector search.
- Alternative C increases index complexity.
- Alternative D makes the index native but adds registry maintenance work.

Consequence: B should generate an internal index from Markdown metadata. The index can power list, diff, validation, and selector commands without making users maintain a separate registry.

## Obsidian Link Analysis

`[[basic/identity#agent-builder]]` is a good human-facing reference shape because it is readable, compact, and familiar to Obsidian users.

Issues:

- Standard Markdown does not define wikilinks.
- Obsidian normally treats the part after `#` as a heading, not an arbitrary metadata ID.
- Headings can be duplicated, renamed, or reworded.
- Paths can move, so presets need either migration support or validation errors.

Surviving interpretation:

- Use an Obsidian-style reference string for presets and selector output.
- Define mogent's semantics explicitly: the `#agent-builder` segment resolves to `mogent.id`, not necessarily the displayed heading text.
- Keep TOML responsible for render category order.

## File Size and Standardization

The project should avoid supporting both file-level and heading-level selection as equal first-class standards in the first redesign. Supporting both is attractive during migration, but it complicates validation, diff output, selector behavior, and user explanation.

Recommended standard:

- Files are containers.
- Blocks are selectable units.
- Every selectable block has an explicit stable ID in metadata.
- A whole-file include can be represented by a file containing one top-level block, rather than a second selection model.

## Conclusions

Rejected alternatives:

- Alternative A as the long-term standard, because file-level modules become too granular or force unrelated options into the same selectable unit.
- Alternative C as the long-term standard, because dual semantics create avoidable complexity.
- Content hashes as primary block IDs, because normal prose edits would break presets.
- TOML-only metadata, because selector descriptions and conflict groups should travel with the block text.

Surviving alternatives:

- Alternative B: heading-level selectable blocks in category files.
- Alternative D as a later generated-index implementation detail, not as the authoring surface.

Recommended conclusion:

- Standardize on heading-level selectable blocks.
- Use TOML only for render order and high-level config.
- Use Obsidian-style references in presets, with mogent-defined semantics.
- Resolve `[[path#id]]` against explicit block metadata IDs.
- Use HTML comments with YAML-like `mogent:` metadata rather than `%%` comments for wider Markdown compatibility.
- Defer tags until block IDs, validation, presets, and diff output are stable.

## Implications for Open TODOs and Pending DIs

The existing `DI-jusuk` decision remains valid for the first implementation, but a new DI is needed before code changes because the recommended block model changes the selection unit and diff model.

Required DF questions before implementation:

1. Selectable unit:
   - A. heading-level blocks only
   - B. file-level modules only
   - C. support both during migration

2. Preset reference syntax:
   - A. Obsidian-style `[[path#id]]`
   - B. plain strings like `path#id`
   - C. TOML objects with `path` and `id`

3. Block ID policy:
   - A. explicit human-authored `mogent.id`
   - B. heading-derived anchors
   - C. content hashes

4. Metadata format:
   - A. HTML comments with YAML-like `mogent:` metadata
   - B. Obsidian `%% ... %%` comments
   - C. TOML-only metadata

5. Render policy:
   - A. category and block order from TOML
   - B. filesystem order
   - C. preset order

6. Diff policy:
   - A. compare selected block IDs plus rendered text
   - B. compare rendered text only
   - C. compare metadata/index only

Decision status: needs DF.
