package commands

type Command interface {
	Run(left, right string) Result
	Description() string
}

type Result struct {
	Affected int
	Success  bool
	Error    error
}
