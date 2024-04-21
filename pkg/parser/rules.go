package parser

type Rule struct {
    If string
    When string
    AllowFailure bool
    AllowFailureExitCodes []int
    Changes []string
}

func (rule *Rule) Parse(template any) error {
	return nil
}
