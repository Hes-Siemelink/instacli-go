package testing

type TestCaseCommand struct{}

func NewTestCaseCommand() *TestCaseCommand {
	return &TestCaseCommand{}
}

func (c *TestCaseCommand) Execute() error {
	// No-op: used as a marker for test boundaries
	return nil
}
