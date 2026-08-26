package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	inputs []textinput.Model
	focus int
	result string
}

func initialModel() Model {
	var inputs []textinput.Model

	c2c := textinput.New()
	c2c.Placeholder = "Target distance"
	c2c.Prompt = "Target C2C Distance (in): "
	c2c.Focus()
	inputs = append(inputs, c2c)

	ratio := textinput.New()
	ratio.Placeholder = "1.0"
	ratio.Prompt = "Target Ratio: "
	inputs = append(inputs, ratio)

	return  Model{
		inputs: inputs,
		focus: 0,
		result: "Waiting on input...",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String(){
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		
		case "up", "down", "tab", "shift+tab":
			if msg.String() == "up" || msg.String() == "shift+tab"{
				m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				m.focus = (m.focus + 1) % len(m.inputs)
			}

			for i := range m.inputs {
				if i == m.focus {
					cmd := m.inputs[i].Focus()
					cmds = append(cmds, cmd)
				} else {
					m.inputs[i].Blur()
				}
			}
		}
	}

	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	c2cVal := m.inputs[0].Value()
	ratioVal := m.inputs[1].Value()
	m.result = fmt.Sprintf("Calcualting for C2C: %s, Ratio %s", c2cVal, ratioVal)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return cmd
}

func (m Model) View() string {
	s := "Pulley Optimizer\n\n"
	for i := range m.inputs {
		s += m.inputs[i].View() + "\n"
	}
	s += "\nTop Results\n"
	s += m.result
	s += "\n\nPress `up/down` to navigate, `q` to quit"
	return  s
}

func main(){
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}