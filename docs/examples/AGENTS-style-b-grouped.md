# Style B: Grouped by Category

# Modules grouped under category headers.
# Each category can have default tags that apply to all its modules.
# Module-specific tags override or add to category tags.

[project]
name = "my-project"
org = "acme"

# Base rules - always included, no special scope
[category.base]
tags = []  # base has no scope tags, always included

[[category.base.module]]
name = "base-rules"
source = "~/.agent-modules/base-rules.md"

[[category.base.module]]
name = "comment-protocol"
source = "~/.agent-modules/comment-protocol.md"

# Language-specific rules
[category.lang]
tags = ["org:acme"]  # all lang modules inherit org:acme

[[category.lang.module]]
name = "go-style"
source = "~/.agent-modules/go-style.md"
tags = ["lang:go"]  # adds lang:go on top of org:acme

[[category.lang.module]]
name = "python-style"
source = "~/.agent-modules/python-style.md"
tags = ["lang:python"]

# Topic-specific rules
[category.topic]
tags = ["org:acme"]

[[category.topic.module]]
name = "security"
source = "~/.agent-modules/security.md"
tags = ["topic:security"]

# Personal overrides - only included for this person
[category.personal]
tags = ["person:quincy"]

[[category.personal.module]]
name = "my-preferences"
source = "~/my-agent-modules/preferences.md"

[[category.personal.module]]
name = "my-go-overrides"
source = "~/my-agent-modules/go-overrides.md"
tags = ["lang:go"]

# Output
[output]
path = "AGENTS.md"
