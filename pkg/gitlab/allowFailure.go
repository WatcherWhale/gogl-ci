package gitlab

type AllowFailure struct {
	AllowFailure          bool  `default:"false"`
	AllowFailureExitCodes []int `default:"[]"`
}

func (a *AllowFailure) Parse(template any) error {
	return nil
}
