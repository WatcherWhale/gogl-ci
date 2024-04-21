package parser

import "github.com/creasty/defaults"

type Artifacts struct {
	Paths    []string `defaults:"[]"`
	Exclude  []string `defaults:"[]"`
	ExpireIn string
	ExposeAs string
	Name     string
	Public   bool
	When     string
}

func (a *Artifacts) Parse(template any) error {
	err := defaults.Set(a)
	if err != nil {
		return err
	}

	return nil
}
