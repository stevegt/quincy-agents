# Style C: Scope-First Organization

# Organized by scope level. Each scope declares what modules it contributes.
# When assembling, active scopes are unioned together.
# If a module appears in multiple scopes, the most specific scope wins.

[project]
name = "my-project"

# Global scope - always active, baseline rules
[scope.global]
source = "~/.agent-modules/"  # default module directory

[[scope.global.module]]
name = "base-rules"
file = "base-rules.md"

[[scope.global.module]]
name = "comment-protocol"
file = "comment-protocol.md"

# Org scope - activated when org matches
[scope.org.acme]
source = "~/.agent-modules/orgs/acme/"

[[scope.org.acme.module]]
name = "go-style"
file = "go-style.md"

[[scope.org.acme.module]]
name = "security"
file = "security.md"

# Team scope - activated when team matches
[scope.team.backend]
source = "~/.agent-modules/teams/backend/"

[[scope.team.backend.module]]
name = "backend-conventions"
file = "conventions.md"

# Person scope - always active for this person
[scope.person.quincy]
source = "~/my-agent-modules/"

[[scope.person.quincy.module]]
name = "preferences"
file = "preferences.md"

[[scope.person.quincy.module]]
name = "go-overrides"
file = "go-overrides.md"
override = "org.acme.go-style"  # replaces the org version

# Active scopes for this repo (resolved at build time)
[activate]
org = "acme"
team = "backend"
person = "quincy"

# Output
[output]
path = "AGENTS.md"
