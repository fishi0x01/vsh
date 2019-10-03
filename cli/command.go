package cli

// Command interface to describe a command structure
type Command interface {
	Run() error
	GetName() string
}
