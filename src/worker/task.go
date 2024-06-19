package worker

type Task interface {
	Execute(worker int)
}
