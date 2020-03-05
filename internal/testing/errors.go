package testing

type TestError string

var _ error = TestError("")

func NewTestError(s string) error {
	return TestError(s)
}

func (e TestError) Error() string {
	return string(e)
}
