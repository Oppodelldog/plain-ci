package build

type Status int

const (
	Queued   Status = 1
	Building Status = 2
	Finished Status = 3
)

func (status Status) String() string {
	if status < Queued || status > Finished {
		status = 0
	}

	return []string{
		"Unknown",
		"Queued",
		"Building",
		"Finished",
	}[status]
}

type Result int

const (
	Success Result = 1
	Failure Result = 2
	Aborted Result = 3
)

func (result Result) String() string {

	if result < Success || result > Aborted {
		result = 0
	}

	return []string{
		"Unknown",
		"Success",
		"Failure",
		"Aborted",
	}[result]
}
