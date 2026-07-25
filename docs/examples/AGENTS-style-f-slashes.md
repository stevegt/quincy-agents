# Style F: Grouped with Slash Hierarchies

# Tags use slash notation: org/acme/team/backend
# Like file paths - intuitive and clean.

[project]
name = "my-project"

# Base rules - always included
[category.base]

[[category.base.module]]
name = "base-rules"
source = "~/.agent-modules/base-rules.md"

# Language rules
[category.lang]
tags = ["org/acme"]  # inherit org

[[category.lang.module]]
name = "go-style"
source = "~/.agent-modules/go-style.md"
tags = ["lang/go"]

[[category.lang.module]]
name = "python-style"
source = "~/.agent-modules/python-style.md"
tags = ["lang/python"]

# Team-specific
[category.team]
tags = ["org/acme"]

[[category.team.module]]
name = "backend-conventions"
source = "~/.agent-modules/teams/backend.md"
tags = ["org/acme/team/backend"]

# Personal
[category.personal]
tags = ["person/quincy"]

[[category.personal.module]]
name = "my-preferences"
source = "~/my-agent-modules/preferences.md"

# Active scopes
[activate]
scopes = [
    "org/acme",
    "org/acme/team/backend",
    "person/quincy",
    "lang/go",
]

# Hierarchical matching:
# - org/acme/team/backend implies org/acme
# - Module tagged ["org/acme"] included when org/acme/team/backend active
# - Module tagged ["org/acme/team/backend"] NOT included when only org/acme active

[output]
path = "AGENTS.md"
