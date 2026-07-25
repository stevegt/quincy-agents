# Style D: Hybrid - Flat with Inheritance

# Flat module list, but each module can declare its scope and inheritance.
# Simple to read, flexible when needed.

[project]
name = "my-project"

# Module sources - can be local, remote, or promise-grid CID
[source.local] = "~/.agent-modules/"
[source.remote] = "git@github.com:acme/agent-modules.git//"
# [source.grid] = "cid://..."  # future promise-grid support

# Modules in assembly order
[[module]]
name = "base-rules"
source = "local:base-rules.md"
# no tags = always included

[[module]]
name = "go-style"
source = "local:go-style.md"
tags = ["lang:go", "org:acme"]

[[module]]
name = "go-overrides"
source = "person:go-overrides.md"
tags = ["lang:go", "person:quincy"]
override = "go-style"  # replaces go-style when both active

[[module]]
name = "security"
source = "remote:security.md"
tags = ["topic:security", "org:acme"]

[[module]]
name = "my-preferences"
source = "person:preferences.md"
tags = ["person:quincy"]

# Inline scope tags in module files override TOML tags
# e.g., a section in go-style.md tagged {#person:quincy} only appears
# when person:quincy is active

# Active scopes
[activate]
org = "acme"
team = "backend"
person = "quincy"

# Output
[output]
path = "AGENTS.md"
