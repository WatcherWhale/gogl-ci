package gitlab

type Artifacts struct {
	Paths    []string `default:"[]"`
	Exclude  []string `default:"[]"`
	ExpireIn string
	ExposeAs string
	Name     string
	Public   bool
	When     string
}

func (a *Artifacts) Parse(template any) error {
	return nil
}
