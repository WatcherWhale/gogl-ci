package format

type TestOutput struct {
	Name string

	Succeeded bool
	Message   string

	SubTests []TestOutput
}

func (t TestOutput) IsTreeSucceeded() bool {
	if !t.Succeeded {
		return false
	}

	for _, st := range t.SubTests {
		if !st.IsTreeSucceeded() {
			return false
		}
	}

	return true
}
