package check

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Check struct {
	Result bool
	Error  error
}

const (
	GoodState = "✔︎"
	BadState  = "⨯"
	Bullet    = "•"
)

var (
	BoldWhite     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	BadStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	TraceBadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingLeft(4)
	WarningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).PaddingLeft(4)
	GoodStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
)

func Checker(check string, err error) {

	if err != nil {
		fmt.Println(BadStyle.Render(BadState + " - " + check))
		fmt.Println(TraceBadStyle.Render(Bullet + " Error: " + err.Error()))

	} else {
		fmt.Println(GoodStyle.Render(GoodState + " - " + check))
	}
}
