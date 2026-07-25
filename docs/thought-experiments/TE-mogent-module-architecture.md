# TE: Mogent Module Architecture

## Decision Under Test
How should mogent organize modules, handle tags, and support different workflows?

## Alternatives

### Alternative A: Flat Structure with Tags
```
.mogent/
  identity.md
  instructions.md
  constraints.md
  format.md
  decision-protocol.md
  coding-style.md
```

**Pros:**
- Simple to understand
- Easy to see all files
- No nested directories

**Cons:**
- Can get cluttered with many files
- Hard to see which files belong together

### Alternative B: Categories as Folders
```
.mogent/
  identity/
    project-overview.md
    team-info.md
  instructions/
    workflow.md
    decision-protocol.md
  constraints/
    security.md
    error-handling.md
  format/
    coding-style.md
    diff-discipline.md
```

**Pros:**
- Clear organization
- Easy to see what's in each category
- Scales well

**Cons:**
- More complex to navigate
- Harder to find a specific file

### Alternative C: Hybrid - Categories for Organization, Tags for Filtering
```
.mogent/
  identity/
    project-overview.md {#team/backend}
    team-info.md {#team/frontend}
  instructions/
    workflow.md
    decision-protocol.md {#process/decisions}
  constraints/
    security.md {#org/acme}
    error-handling.md
  format/
    coding-style.md {#lang/go}
    diff-discipline.md
```

**Pros:**
- Clear organization by category
- Flexible filtering by tags
- Supports multiple workflows

**Cons:**
- More complex to set up
- Requires understanding both categories and tags

## Workflow Scenarios

### Scenario 1: Solo Developer, Simple Project
- Just 4 files: identity, instructions, constraints, format
- No tags needed
- Simple `mogent build` to assemble

### Scenario 2: Team with Multiple Projects
- Tags for team: team/backend, team/frontend
- Tags for project: project/web, project/mobile
- `mogent build --tags team/backend` to get backend-specific AGENTS.md

### Scenario 3: Large Org with Many Repos
- Tags for org: org/acme
- Tags for repo: repo/agents, repo/mint
- `mogent build --tags org/acme,repo/agents` for specific repo

### Scenario 4: Multiple Workflows
- Tags for workflow: workflow/prototyping, workflow/production
- Different files for different stages
- `mogent build --tags workflow/production` for production rules

## Design Questions

### 1. Where do TE/DI/DF protocols go?
- Option A: In `instructions/decision-protocol.md`
- Option B: In `constraints/decision-protocol.md`
- Option C: Separate `process/` category

### 2. Should categories be folders?
- Option A: Yes, always
- Option B: No, always flat
- Option C: Configurable in AGENTS.toml

### 3. How to pull files by tag vs by file?
- Option A: `mogent build --tags team/backend`
- Option B: `mogent build --files identity,instructions`
- Option C: Both

### 4. What templates should we have?
- Option A: Just the 4 basic ones
- Option B: Add decision-protocol, coding-style, error-handling
- Option C: Full set with all protocols

## Recommendations

1. **Use Alternative C (Hybrid)** - Categories for organization, tags for filtering
2. **TE/DI/DF go in `instructions/decision-protocol.md`** - It's about how to work, not what's forbidden
3. **Categories should be folders** - Scales better, clearer organization
4. **Support both tag and file selection** - `mogent build --tags team/backend` and `mogent build --files identity,instructions`
5. **Start with basic templates** - Add more as needed

## Next Steps

1. Implement tag management commands
2. Add list command to show all modules and tags
3. Create example files for each category
4. Test with different workflow scenarios
