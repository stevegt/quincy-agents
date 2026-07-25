# Style A: Flat List with Category Tags

# Simple flat list. Category is just another tag.
# Assembly order is TOML order.
# Each module can have inline scope tags.

[project]
name = "my-project"
org = "acme"

[[module]]
name = "base-rules"
source = "~/.agent-modules/base-rules.md"
tags = ["base"]

[[module]]
name = "go-style"
source = "~/.agent-modules/go-style.md"
tags = ["lang:go", "org:acme"]

[[module]]
name = "security"
source = "~/.agent-modules/security.md"
tags = ["topic:security", "org:acme"]

[[module]]
name = "my-additions"
source = "~/my-agent-modules/custom.md"
tags = ["person:quincy"]

# Output
[output]
path = "AGENTS.md"
