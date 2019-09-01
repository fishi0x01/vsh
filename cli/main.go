package cli

type Command interface {
	Run() error
	GetName() string
}
