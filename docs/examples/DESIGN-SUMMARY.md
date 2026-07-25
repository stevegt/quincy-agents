# Mogent Design Summary

## Locked Decisions

| Decision | Choice | Notes |
|----------|--------|-------|
| CLI Language | Go | Matches wire-lab ecosystem |
| CLI Name | mogent | Working title, modular agent |
| Scope Model | Tag-based | Inline tags in module files |
| Tag Separator | Slashes | org/acme/team/backend |
| Hierarchical Matching | Parent implies children | org/acme/team/backend implies org/acme |
| Inline Tag Syntax | Header attributes | ## Section {#tag1 #tag2} |
| Module Identity | Name-based | Names in TOML, paths for resolution |
| Assembly Order | Category groups | TOML order within categories |
| Diff Model | Section-aware + semantic | Parse markdown into blocks |
| Promise Grid | Design for later | Keep architecture open |
| Module Storage | .mogent/ default | Global default tag path, per-module overrides |
| Editor Integration | Both modes | Edit module directly, or edit assembled output |
| Conflict Resolution | Most specific wins | person/quincy overrides org/acme |

## AGENTS.toml Structure (Final)

```toml
# Default module directory
[config]
module_dir = ".mogent/modules"  # default, can be overridden per category

# Base rules - always included
[category.base]

[[category.base.module]]
name = "base-rules"
source = "base-rules.md"  # resolves to .mogent/modules/base-rules.md

# Language rules - inherit org
[category.lang]
tags = ["org/acme"]

[[category.lang.module]]
name = "go-style"
source = "go-style.md"
tags = ["lang/go"]

# Team-specific
[category.team]
tags = ["org/acme"]

[[category.team.module]]
name = "backend-conventions"
source = "teams/backend.md"
tags = ["org/acme/team/backend"]

# Personal overrides
[category.personal]
tags = ["person/quincy"]

[[category.personal.module]]
name = "my-preferences"
source = "~/my-mogent/preferences.md"  # absolute path override

# Active scopes for this repo
[activate]
scopes = [
    "org/acme",
    "org/acme/team/backend",
    "person/quincy",
    "lang/go",
]

# Output
[output]
path = "AGENTS.md"
```

## Module File Format

```markdown
# Go Style Guide

## General Rules {#org/acme #lang/go}
All Go code must use gofmt.

## Error Handling {#org/acme #lang/go}
Use %w for error wrapping.

## Personal Style {#person/quincy}
I prefer early returns over nested if.
```

## CLI Commands

```
mogent build          # Assemble AGENTS.md from modules
mogent build --dry    # Preview without writing
mogent edit <module>  # Edit module in $EDITOR
mogent edit           # Edit assembled AGENTS.md
mogent diff           # Compare local vs assembled
mogent diff <scope>   # Compare with specific scope
mogent list           # List available modules
mogent list --active  # List modules included for active scopes
```

## Directory Structure

```
repo/
├── AGENTS.toml              # Project config
├── AGENTS.md                # Assembled output (generated)
├── .mogent/
│   └── modules/
│       ├── base-rules.md
│       ├── go-style.md
│       ├── security.md
│       └── teams/
│           └── backend.md
└── ...

~/.mogent/                   # User-global modules
├── modules/
│   ├── base-rules.md
│   ├── go-style.md
│   └── ...
└── config.toml              # Global defaults
```

## Promise Grid Integration Points

1. **Module as CAS object**: Each module file can be addressed by CID
2. **Remote modules**: Source can be a CID reference instead of file path
3. **Version tracking**: Module changes tracked as reference sets
4. **Cross-org sharing**: Share modules via grid sync, not file copy

Architecture keeps a `Source` interface that can be:
- Local file path
- Remote URL
- CID reference (future)
- Grid sync reference (future)
