package handler

import (
	"gostatus/internal/store"
)

// FindActivity returns the first activity matching the given type and name, or nil.
func FindActivity(p store.Presence, actType int, name string) *store.Activity {
	for i := range p.Activities {
		if a := p.Activities[i]; a.Type == actType && (name == "" || a.Name == name) {
			return &p.Activities[i]
		}
	}
	return nil
}

// FindAllActivities returns all activities of the given type, excluding the given names.
func FindAllActivities(p store.Presence, actType int, exclude ...string) []store.Activity {
	excluded := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		excluded[e] = true
	}
	var out []store.Activity
	for _, a := range p.Activities {
		if a.Type == actType && !excluded[a.Name] {
			out = append(out, a)
		}
	}
	return out
}
