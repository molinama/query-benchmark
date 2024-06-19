package repository

type DBMock struct{}

func (db DBMock) Execute(worker int) {
	panic("not implemented") // TODO: Implement
}
