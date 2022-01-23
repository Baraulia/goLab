package myErrors

type Error interface {
	error
	Status() int
}

type MyError struct {
	Err  error
	Code int
}

func (m *MyError) Error() string {
	return m.Err.Error()
}

func (m *MyError) Status() int {
	return m.Code
}
