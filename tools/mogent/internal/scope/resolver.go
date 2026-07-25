// Intent: Hierarchical tag matching with parent-implies-children semantics.
// Source: DI-jusuk

package scope

import "strings"

type Resolver struct {
	active map[string]bool
}

func NewResolver(activeScopes []string) *Resolver {
	r := &Resolver{
		active: make(map[string]bool),
	}

	for _, scope := range activeScopes {
		r.active[scope] = true
		parts := strings.Split(scope, "/")
		for i := 1; i < len(parts); i++ {
			r.active[strings.Join(parts[:i], "/")] = true
		}
	}

	return r
}

func (r *Resolver) Matches(tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	for _, tag := range tags {
		if r.active[tag] {
			return true
		}
	}

	return false
}

func (r *Resolver) Specificity(tags []string) int {
	maxSpec := 0
	for _, tag := range tags {
		if r.active[tag] {
			parts := strings.Split(tag, "/")
			if len(parts) > maxSpec {
				maxSpec = len(parts)
			}
		}
	}
	return maxSpec
}
