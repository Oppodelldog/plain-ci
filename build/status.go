package build

type BuildStatus int

const (
	Queued   BuildStatus = 0
	Building BuildStatus = 1
	Finished BuildStatus = 2
)

func (status BuildStatus) String() string {
	names := [...]string{
		"Queued",
		"Building",
		"Finished",
	}

	if status < Queued || status > Finished {
		return "Unknown"
	}

	return names[status]
}
