package handler

import (
	"gostatus/internal/gateway"
	"gostatus/internal/store"
)

type Activity = gateway.Activity

func FindActivity(p store.Presence, actType int, name string) *Activity {
	for _, a := range p.Activities {
		if a.Type == actType && (name == "" || a.Name == name) {
			a := Activity(a)
			return &a
		}
	}
	return nil
}

func FindAllActivities(p store.Presence, actType int, exclude ...string) []Activity {
	excluded := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		excluded[e] = true
	}
	var out []Activity
	for _, a := range p.Activities {
		if a.Type == actType && !excluded[a.Name] {
			out = append(out, Activity(a))
		}
	}
	return out
}
