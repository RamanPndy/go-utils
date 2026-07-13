package goutils

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ChangeType classifies a single field change.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// FieldChange records a change to a single JSON path.
type FieldChange struct {
	Path     string     `json:"path"`
	Type     ChangeType `json:"type"`
	OldValue any        `json:"old_value,omitempty"`
	NewValue any        `json:"new_value,omitempty"`
}

// Diff is the collection of all field changes between two objects.
type Diff struct {
	Changes []FieldChange `json:"changes"`
}

// Empty returns true when no changes were found.
func (d *Diff) Empty() bool {
	return len(d.Changes) == 0
}

// Has returns true when the diff contains a change at the given path prefix.
func (d *Diff) Has(pathPrefix string) bool {
	for _, c := range d.Changes {
		if len(c.Path) >= len(pathPrefix) && c.Path[:len(pathPrefix)] == pathPrefix {
			return true
		}
	}
	return false
}

// Get returns the FieldChange for an exact path, or nil.
func (d *Diff) Get(path string) *FieldChange {
	for i := range d.Changes {
		if d.Changes[i].Path == path {
			return &d.Changes[i]
		}
	}
	return nil
}

// Compute calculates the diff between oldRaw and newRaw JSON bytes.
// Both inputs should be valid JSON objects; if either is nil the function
// treats it as an empty object.
func Compute(oldRaw, newRaw []byte) (*Diff, error) {
	oldMap, err := toMap(oldRaw)
	if err != nil {
		return nil, fmt.Errorf("parsing old object: %w", err)
	}
	newMap, err := toMap(newRaw)
	if err != nil {
		return nil, fmt.Errorf("parsing new object: %w", err)
	}

	d := &Diff{}
	walkDiff("", oldMap, newMap, d)
	sortChanges(d)
	return d, nil
}

// ToJSON serialises the diff as JSON.
func (d *Diff) ToJSON() ([]byte, error) {
	return json.Marshal(d.Changes)
}

// walkDiff recursively walks two maps comparing their values.
func walkDiff(prefix string, old, new map[string]any, d *Diff) {
	// Keys in old
	for k, oldVal := range old {
		path := joinPath(prefix, k)
		newVal, exists := new[k]
		if !exists {
			d.Changes = append(d.Changes, FieldChange{Path: path, Type: ChangeRemoved, OldValue: oldVal})
			continue
		}
		compareValues(path, oldVal, newVal, d)
	}

	// Keys only in new
	for k, newVal := range new {
		path := joinPath(prefix, k)
		if _, exists := old[k]; !exists {
			d.Changes = append(d.Changes, FieldChange{Path: path, Type: ChangeAdded, NewValue: newVal})
		}
	}
}

// compareValues compares two arbitrary JSON values at the given path.
func compareValues(path string, oldVal, newVal any, d *Diff) {
	oldMap, oldIsMap := oldVal.(map[string]any)
	newMap, newIsMap := newVal.(map[string]any)

	if oldIsMap && newIsMap {
		walkDiff(path, oldMap, newMap, d)
		return
	}

	if !jsonEqual(oldVal, newVal) {
		d.Changes = append(d.Changes, FieldChange{
			Path:     path,
			Type:     ChangeModified,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}
}

// toMap unmarshals raw JSON into a map; nil input returns an empty map.
func toMap(raw []byte) (map[string]any, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// jsonEqual compares two JSON-unmarshalled values for equality.
func jsonEqual(a, b any) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}

// joinPath builds a dot-separated JSON path.
func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// sortChanges ensures deterministic ordering of diff output.
func sortChanges(d *Diff) {
	sort.Slice(d.Changes, func(i, j int) bool {
		return d.Changes[i].Path < d.Changes[j].Path
	})
}
