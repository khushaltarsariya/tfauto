package cmd

type ExitError struct {
	Code    int
	Message string
}

func (e ExitError) Error() string {
	return e.Message
}

func (e ExitError) ExitCode() int {
	if e.Code == 0 {
		return 1
	}

	return e.Code
}
