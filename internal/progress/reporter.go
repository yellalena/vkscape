package progress

type Reporter interface {
	Start(total int)
	Increment()
	SetStatus(msg string)
	Done()
}

type NoopReporter struct{}

func (*NoopReporter) Start(int)        {}
func (*NoopReporter) Increment()       {}
func (*NoopReporter) SetStatus(string) {}
func (*NoopReporter) Done()            {}
