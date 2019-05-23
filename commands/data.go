package commands

type Command interface {
	Run(left, right string) Result
}

type Result struct {
	Affected int
	Success  bool
	Error    error
}
