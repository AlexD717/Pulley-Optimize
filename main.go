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
	TypeText     FieldType = "text"
	TypeCheckbox FieldType = "checkbox"
)

type FormField struct {
	Name    string
	Type    FieldType
	Input   textinput.Model
	Checked bool // Used only if Checkbox
	Visible bool
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
	c2c.Prompt = "Target C2C Distance (in): "
	c2c.Validate = numberValidator
	c2c.Focus()

	ratio := textinput.New()
	ratio.Placeholder = "1"
	ratio.Prompt = "Target Ratio: "
	ratio.Validate = numberValidator

	fields := []FormField{
		{Name: "C2C", Type: TypeNumber, Input: c2c, Visible: true},
		{Name: "Ratio", Type: TypeNumber, Input: ratio, Visible: true},
		{Name: "Use Available Belts", Type: TypeCheckbox, Visible: true},
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
				m.focus = (m.focus - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				m.focus = (m.focus + 1) % len(m.inputs)
			}

			var cmd tea.Cmd
			m.inputs, cmd = updateFocusStyles(m.inputs, m.focus)
			cmds = append(cmds, cmd)

		case "left", "right":
			passToInput = false
			valStr := m.inputs[m.focus].Input.Value()
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				val = 0
			}

			step := 0.1
			if m.focus == 1 {
				step = .2
			}

			if msg.String() == "right" {
				val += step
			} else if msg.String() == "left" {
				val -= step
			}

			m.inputs[m.focus].Input.SetValue(fmt.Sprintf("%.1f", val))

		case " ", "enter":
			if m.inputs[m.focus].Type == TypeCheckbox {
				m.inputs[m.focus].Checked = !m.inputs[m.focus].Checked
			}
		}
	}

	if passToInput {
		cmd := m.updateInputs(msg)
		cmds = append(cmds, cmd)
	}

	c2cVal := m.inputs[0].Input.Value()
	ratioVal := m.inputs[1].Input.Value()
	useNotion := m.inputs[2].Checked
	m.result = fmt.Sprintf("Calculating for C2C: %s, Ratio %s, Using Notion %v", c2cVal, ratioVal, useNotion)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.inputs[m.focus].Type != TypeCheckbox {
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

		if f.Type == TypeCheckbox {
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
		} else {
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
		if fields[i].Type == TypeCheckbox {
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
