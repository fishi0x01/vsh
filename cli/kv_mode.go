package cli

// KeyValueMode defines the scope of which parts of a path to search (keys and/or values)
type KeyValueMode int

const (
	// ModeKeys only searches keys
	ModeKeys KeyValueMode = 1
	// ModeValues only searches values
	ModeValues KeyValueMode = 2
)

// KeyValueCommand interface to describe a command that supports Key and/or Value scoping
type KeyValueCommand interface {
	IsMode(mode KeyValueMode) bool
}
