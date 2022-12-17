package types

type NotSupportCommandError struct{}

func (nsc NotSupportCommandError) Error() string {
	return "not support command usage"
}
