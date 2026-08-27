package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type FieldType string

const (
	TypeNumber   FieldType = "number"
	TypeCheckbox FieldType = "checkbox"
	TypeSelector FieldType = "selector"
)

type FormField struct {
	Name     string
	Type     FieldType
	Input    textinput.Model
	Checked  bool // Used only if checkbox
	Visible  bool
	Options  []string // Used by selector
	Selected int      // Used by selector
}

type Model struct {
	Inputs            []FormField
	Focus             int
	UseAvailableBelts bool
	Results           []PulleyResult
	ErrorMsg          string
}

func initialModel() Model {
	c2c := textinput.New()
	c2c.Placeholder = "0"
	c2c.Prompt = "Target C2C Distance: "
	c2c.Focus()

	ratio := textinput.New()
	ratio.Placeholder = "1.0"
	ratio.SetValue("1.0")
	ratio.Prompt = "Target Ratio: "

	fields := []FormField{
		{Name: "C2C", Type: TypeNumber, Input: c2c, Visible: true},
		{Name: "Unit", Type: TypeSelector, Visible: true, Options: []string{"in", "mm"}, Selected: 0},
		{Name: "Ratio", Type: TypeNumber, Input: ratio, Visible: true},
		{Name: "Use Available Belts", Type: TypeCheckbox, Visible: true, Checked: true},
	}

	fields, _ = updateFocusStyles(fields, 0)

	return Model{
		Inputs:            fields,
		Focus:             0,
		UseAvailableBelts: true,
		Results:           []PulleyResult{},
		ErrorMsg:          "",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	passToInput := true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "down", "tab", "shift+tab":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				for {
					m.Focus = (m.Focus - 1 + len(m.Inputs)) % len(m.Inputs)
					if m.Inputs[m.Focus].Visible {
						break
					}
				}
			} else {
				for {
					m.Focus = (m.Focus + 1) % len(m.Inputs)
					if m.Inputs[m.Focus].Visible {
						break
					}
				}
			}

			var cmd tea.Cmd
			m.Inputs, cmd = updateFocusStyles(m.Inputs, m.Focus)
			cmds = append(cmds, cmd)

		case "left", "right":
			passToInput = false
			f := &m.Inputs[m.Focus]
			switch f.Type {
			case TypeNumber:
				valStr := m.Inputs[m.Focus].Input.Value()
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					val = 0
				}

				step := 0.1
				if m.Inputs[m.Focus].Name == "Ratio" {
					step = .2
				}

				if msg.String() == "right" {
					val += step
				} else if msg.String() == "left" {
					val -= step
				}

				f.Input.SetValue(fmt.Sprintf("%.1f", val))
			case TypeSelector:
				if msg.String() == "right" {
					f.Selected = (f.Selected + 1) % len(f.Options)
				} else if msg.String() == "left" {
					f.Selected = (f.Selected - 1 + len(f.Options)) % len(f.Options)
				}
			case TypeCheckbox:
				f.Checked = !f.Checked
			}

		case " ", "enter":
			f := &m.Inputs[m.Focus]
			switch f.Type {
			case TypeCheckbox:
				f.Checked = !f.Checked
			case TypeSelector:
				f.Selected = (f.Selected + 1) % len(f.Options)
			}
		}
	}

	if passToInput {
		cmd := m.updateInputs(msg)
		cmds = append(cmds, cmd)
	}

	m = m.updateResults()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.Inputs[m.Focus].Type == TypeNumber {
		m.Inputs[m.Focus].Input, cmd = m.Inputs[m.Focus].Input.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	title := titleStyle.Render("FRC Pulley Optimizer ")
	slashes := slashStyle.Render("///////////////////////////////////////////////////////////////////////////////////////////")
	header := title + slashes + "\n\n"

	leftColumn := ""
	for i, f := range m.Inputs {
		if !f.Visible {
			continue
		}

		cursor := ""
		if i == m.Focus {
			cursor = "> "
		}

		switch f.Type {
		case TypeCheckbox:
			box := "[ ]"
			if f.Checked {
				box = "[x]"
			}

			styledName := unActiveStyle.Render(f.Name)
			if i == m.Focus {
				styledName = activeStyle.Render(f.Name)
			}
			leftColumn += fmt.Sprintf("%s %s\n", cursor+styledName, box)
		case TypeSelector:
			styledName := unActiveStyle.Render(f.Name)
			if i == m.Focus {
				styledName = activeStyle.Render(f.Name)
			}

			optsStr := ""
			for j, opt := range f.Options {
				if j == f.Selected {
					optsStr += activeTextStyle.Render(fmt.Sprintf("[%s] ", opt))
				} else {
					optsStr += unActiveTextStyle.Render(fmt.Sprintf(" %s  ", opt))
				}
			}
			leftColumn += fmt.Sprintf("%s %s\n", cursor+styledName, optsStr)

		default:
			inputView := f.Input.View()

			if f.Input.Err != nil {
				inputView += errorStyle.Render("  <- " + f.Input.Err.Error())
			}

			leftColumn += fmt.Sprintf("%s\n", fixedWidthText.Render(cursor+inputView))
		}
	}

	boxTitle := boxTitleStyle.Render("Top Results")

	var resultsString string

	if m.ErrorMsg != "" {
		resultsString = errorStyle.Render(m.ErrorMsg)
	} else if len(m.Results) == 0 {
		resultsString = "No valid combination found"
	} else {
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(resultsTableStyle).
			Headers("Pulley 1", "Pulley 2", "Ratio", "Belt Length", "Belt Width", "Slack", "Available")

		for _, res := range m.Results {
			var ratio float64 = float64(res.Pulley2) / float64(res.Pulley1)
			t.Row(
				fmt.Sprintf("%dT", res.Pulley1),
				fmt.Sprintf("%dT", res.Pulley2),
				fmt.Sprintf("%.2f", ratio),
				fmt.Sprintf("%d mm", res.BeltLength),
				fmt.Sprintf("%.1f mm", res.BeltWidth),
				fmt.Sprintf("%.2f mm", res.Slack),
				formatAvailabilityBool(res.IsAvailable),
			)
		}

		resultsString = t.Render()
	}

	rightColumnContent := boxTitle + "\n\n" + resultsString

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumnStyle.Render(leftColumn),
		resultBoxStyle.Render(rightColumnContent),
	)

	endSlashes := slashStyle.Render("\n\n////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
	footer := endSlashes + helpStyle.Render("\n\nPress `up/down` to navigate, `left/right` to change values, `q` to quit")

	return header + mainContent + footer
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func updateFocusStyles(fields []FormField, focus int) ([]FormField, tea.Cmd) {
	var cmds []tea.Cmd

	for i := range fields {
		if fields[i].Type == TypeCheckbox || fields[i].Type == TypeSelector {
			continue
		}
		if i == focus {
			cmd := fields[i].Input.Focus()
			cmds = append(cmds, cmd)
			fields[i].Input.PromptStyle = activeStyle
			fields[i].Input.TextStyle = activeTextStyle
		} else {
			fields[i].Input.Blur()
			fields[i].Input.PromptStyle = unActiveStyle
			fields[i].Input.TextStyle = unActiveTextStyle
		}
	}

	return fields, tea.Batch(cmds...)
}

func (m Model) updateResults() Model {
	c2cVal := m.Inputs[0].Input.Value()
	unitVal := m.Inputs[1].Options[m.Inputs[1].Selected]
	ratioVal := m.Inputs[2].Input.Value()
	useAvailableBelts := m.Inputs[3].Checked

	results, err := RunCalculator(c2cVal, ratioVal, unitVal, useAvailableBelts)
	if err != nil {
		m.ErrorMsg = err.Error()
		m.Results = nil
	} else {
		m.ErrorMsg = ""
		m.Results = results
	}

	return m
}

func formatAvailabilityBool(isAvailable bool) string {
	if isAvailable {
		return "Yes"
	} else {
		return "No"
	}
}
