package narwhal_lib

type (
	Executable interface {
		Write(b []byte) error
		CustomRun(outEvent OutputEvent, errEvent OutputEvent) string
		Run() []string
	}

	CommandCreator interface {
		Create(command string, arg ...string) Executable
	}

	OutputEvent = func(string)
)
