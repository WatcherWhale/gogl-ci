package parser

type Artifacts struct {
	Paths    []string
	Exclude  []string
	ExpireIn string
	ExposeAs string
	Name     string
	Public   bool
	When     string
}

func (a *Artifacts) Parse(template any) error {
	return nil
}
