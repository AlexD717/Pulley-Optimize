package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	inputs            []FormField
	focus             int
	useAvailableBelts bool
	result            string
}

func initialModel() Model {
	numberValidator := func(s string) error {
		if s == "" || s == "." {
			return nil
		}

		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("Must be a number")
		}
		return nil
	}

	c2c := textinput.New()
	c2c.Placeholder = "0"
	c2c.Prompt = "Target C2C Distance: "
	c2c.Validate = numberValidator
	c2c.Focus()

	ratio := textinput.New()
	ratio.Placeholder = "1.0"
	ratio.SetValue("1.0")
	ratio.Prompt = "Target Ratio: "
	ratio.Validate = numberValidator

	fields := []FormField{
		{Name: "C2C", Type: TypeNumber, Input: c2c, Visible: true},
		{Name: "Unit", Type: TypeSelector, Visible: true, Options: []string{"in", "mm"}, Selected: 0},
		{Name: "Ratio", Type: TypeNumber, Input: ratio, Visible: true},
		{Name: "Use Available Belts", Type: TypeCheckbox, Visible: true, Checked: true},
	}

	fields, _ = updateFocusStyles(fields, 0)

	return Model{
		inputs: fields,
		focus:  0,
		result: "Waiting on input...",
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
					m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
					if m.inputs[m.focus].Visible {
						break
					}
				}
			} else {
				for {
					m.focus = (m.focus + 1) % len(m.inputs)
					if m.inputs[m.focus].Visible {
						break
					}
				}
			}

			var cmd tea.Cmd
			m.inputs, cmd = updateFocusStyles(m.inputs, m.focus)
			cmds = append(cmds, cmd)

		case "left", "right":
			passToInput = false
			f := &m.inputs[m.focus]
			switch f.Type {
			case TypeNumber:
				valStr := m.inputs[m.focus].Input.Value()
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					val = 0
				}

				step := 0.1
				if m.inputs[m.focus].Name == "Ratio" {
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
			f := &m.inputs[m.focus]
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

	c2cVal := m.inputs[0].Input.Value()
	unitVal := m.inputs[1].Options[m.inputs[1].Selected]
	ratioVal := m.inputs[2].Input.Value()
	useAvailableBelts := m.inputs[3].Checked
	m.result = fmt.Sprintf("Calculating for C2C: %s, Ratio %s, Pitch %s, Using Belts %v", c2cVal, ratioVal, unitVal, useAvailableBelts)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.inputs[m.focus].Type == TypeNumber {
		m.inputs[m.focus].Input, cmd = m.inputs[m.focus].Input.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	title := titleStyle.Render("FRC Pulley Optimizer ")
	slashes := slashStyle.Render("///////////////////////////////////////////////////////////////////////////////////////////")
	header := title + slashes + "\n\n"

	leftColumn := ""
	for i, f := range m.inputs {
		if !f.Visible {
			continue
		}

		switch f.Type {
		case TypeCheckbox:
			box := "[ ]"
			if f.Checked {
				box = "[x]"
			}

			cursor := ""
			styledName := unActiveStyle.Render(f.Name)
			if i == m.focus {
				cursor = "> "
				styledName = activeStyle.Render(f.Name)
			}
			leftColumn += fmt.Sprintf("%s%s %s\n", cursor, styledName, box)
		case TypeSelector:
			cursor := ""
			styledName := unActiveStyle.Render(f.Name)
			if i == m.focus {
				cursor = "> "
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
			leftColumn += fmt.Sprintf("%s%s %s\n", cursor, styledName, optsStr)

		default:
			inputView := f.Input.View()

			if f.Input.Err != nil {
				inputView += errorStyle.Render("  <- " + f.Input.Err.Error())
			}

			leftColumn += fixedWidthText.Render(inputView) + "\n"
		}
	}

	boxTitle := boxTitleStyle.Render("Top Results")
	rightColumnContent := boxTitle + "\n\n" + m.result

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumnStyle.Render(leftColumn),
		resultBoxStyle.Render(rightColumnContent),
	)

	endSlashes := slashStyle.Render("\n\n////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
	footer := endSlashes + helpStyle.Render("\n\nPress `up/down` to navigate, `q` to quit")

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
