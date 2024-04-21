package parser

type AllowFailure struct {
	AllowFailure bool `default:"false"`
	AllowFailureExitCodes []int
}

func (a *AllowFailure) Parse(template any) error {
	return nil
}
