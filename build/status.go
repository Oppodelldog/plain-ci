package build

type Status int

const (
	Queued   Status = 0
	Building Status = 1
	Finished Status = 2
)

func (status Status) String() string {
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
