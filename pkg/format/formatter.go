package format

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

var enumStyle lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginRight(1)

func SprintTests(tests TestOutput) string {
	t := generateSubTree(tests)

	t = t.
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumStyle)

	return fmt.Sprintln(t)
}

func PrintTests(tests TestOutput) {
	fmt.Print(SprintTests(tests))
}

func generateSubTree(test TestOutput) *tree.Tree {
	succeeded := test.IsTreeSucceeded()

	var t *tree.Tree

	if succeeded {
		rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
		t = tree.Root(
			rootStyle.Render(
				fmt.Sprintf(" %s", test.Name),
			),
		)
	} else {
		rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)

		if test.Message != "" {
			t = tree.Root(
				rootStyle.Render(
					fmt.Sprintf(" %s: %s", test.Name, test.Message),
				),
			)
		} else {
			t = tree.Root(
				rootStyle.Render(
					fmt.Sprintf(" %s", test.Name),
				),
			)
		}
	}

	for _, st := range test.SubTests {
		t = t.Child(generateSubTree(st))
	}

	return t
}
