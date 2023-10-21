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
	BoldWhite = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	Bad       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	TraceBad  = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingLeft(4)
	Warning   = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).PaddingLeft(4)
	Good      = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
)

func Checker(check string, err error) {

	if err != nil {
		fmt.Println(Bad.Render(BadState + " - " + check))
		fmt.Println(TraceBad.Render(Bullet + " Error: " + err.Error()))

	} else {
		fmt.Println(Good.Render(GoodState + " - " + check))
	}
}
