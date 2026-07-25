# Style E: Grouped with Hierarchical Tags

# Tags use path-like syntax: org:acme/team:backend
# Activating a parent includes all children unless overridden.

[project]
name = "my-project"

# Base rules - always included, no tags
[category.base]

[[category.base.module]]
name = "base-rules"
source = "~/.agent-modules/base-rules.md"

# Language rules - inherit org:acme, add lang:go
[category.lang]
tags = ["org:acme"]  # all lang modules inherit this

[[category.lang.module]]
name = "go-style"
source = "~/.agent-modules/go-style.md"
tags = ["lang:go"]  # adds lang:go, inherits org:acme

[[category.lang.module]]
name = "python-style"
source = "~/.agent-modules/python-style.md"
tags = ["lang:python"]

# Team-specific overrides
[category.team]
tags = ["org:acme"]  # inherit org

[[category.team.module]]
name = "backend-conventions"
source = "~/.agent-modules/teams/backend.md"
tags = ["org:acme/team:backend"]  # specific to backend team

# Personal - only for this person
[category.personal]
tags = ["person:quincy"]

[[category.personal.module]]
name = "my-preferences"
source = "~/my-agent-modules/preferences.md"

# What's active in this repo
[activate]
scopes = [
    "org:acme",           # org-wide rules
    "org:acme/team:backend",  # backend team specific
    "person:quincy",      # personal overrides
    "lang:go",            # Go-specific rules
]

# Assembly order: base -> lang (filtered) -> team (filtered) -> personal
# Within each category, modules are assembled in TOML order

# Resolution rules:
# 1. base modules: always included (no tags)
# 2. lang modules: included if their tags match any activated scope
#    - go-style has tags ["org:acme", "lang:go"] - both active, included
#    - python-style has tags ["org:acme", "lang:python"] - lang:python not active, excluded
# 3. team modules: same tag matching
# 4. personal modules: same tag matching
# 5. If a module is included multiple times, the most specific scope wins
#    e.g., person:quincy overrides org:acme for same module

# Hierarchical tag matching:
# - org:acme/team:backend implies org:acme
# - So a module tagged ["org:acme"] is included when org:acme/team:backend is active
# - But a module tagged ["org:acme/team:backend"] is NOT included when only org:acme is active

[output]
path = "AGENTS.md"
