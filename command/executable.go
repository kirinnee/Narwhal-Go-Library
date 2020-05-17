package command

type (
	Executable interface {
		Write(b []byte) error
		CustomRun(outEvent OutputEvent, errEvent OutputEvent) string
		Run() []string
	}

	Creator interface {
		Create(command string, arg ...string) Executable
	}

	OutputEvent = func(string)
)
