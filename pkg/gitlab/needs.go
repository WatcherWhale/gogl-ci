package gitlab

type Need struct {
	Job       string
	Artifacts bool
}

func (need *Need) Parse(template any) error {
	return nil
}
