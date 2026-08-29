package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
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
	Advanced bool
	Checked  bool // Used only if checkbox
	Visible  bool
	Options  []string // Used by selector
	Selected int      // Used by selector
}

type CalcResultMessage struct {
	results []PulleyResult
	err     error
}

type Model struct {
	Inputs            []FormField
	Focus             int
	UseAvailableBelts bool
	Results           []PulleyResult
	ErrorMsg          string
	Calculating       bool
	Spinner           spinner.Model
	CancelCalc        context.CancelFunc
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

	maxSlack := textinput.New()
	maxSlack.Placeholder = "0.2"
	maxSlack.SetValue("0.2")
	maxSlack.Prompt = "Max Slack (mm): "

	minSlack := textinput.New()
	minSlack.Placeholder = "-0.5"
	minSlack.SetValue("-0.5")
	minSlack.Prompt = "Min Slack (mm): "

	maxPulley := textinput.New()
	maxPulley.Placeholder = "100"
	maxPulley.SetValue("100")
	maxPulley.Prompt = "Max Pulley (T): "

	minPulley := textinput.New()
	minPulley.Placeholder = "8"
	minPulley.SetValue("8")
	minPulley.Prompt = "Min Pulley (T): "

	fields := []FormField{
		{Name: "C2C", Type: TypeNumber, Input: c2c, Advanced: false, Visible: true},
		{Name: "Unit", Type: TypeSelector, Advanced: false, Visible: true, Options: []string{"in", "mm"}, Selected: 0},
		{Name: "Ratio", Type: TypeNumber, Input: ratio, Advanced: false, Visible: true},
		{Name: "Use Available Belts", Type: TypeCheckbox, Advanced: false, Visible: true, Checked: true},
		{Name: "Max Slack", Type: TypeNumber, Input: maxSlack, Advanced: true, Visible: false},
		{Name: "Min Slack", Type: TypeNumber, Input: minSlack, Advanced: true, Visible: false},
		{Name: "Max Pulley", Type: TypeNumber, Input: maxPulley, Advanced: true, Visible: false},
		{Name: "Min Pulley", Type: TypeNumber, Input: minPulley, Advanced: true, Visible: false},
	}

	fields, _ = updateFocusStyles(fields, 0)

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = spinnerStyle

	return Model{
		Inputs:            fields,
		Focus:             0,
		UseAvailableBelts: true,
		Results:           []PulleyResult{},
		ErrorMsg:          "",
		Calculating:       false,
		Spinner:           s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.Spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	passToInput := true
	inputChanged := false

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)

	case CalcResultMessage:
		if msg.err == context.Canceled {
			return m, nil
		}

		m.Calculating = false
		if msg.err != nil {
			m.ErrorMsg = msg.err.Error()
			m.Results = nil
		} else {
			m.ErrorMsg = ""
			m.Results = msg.results
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "a":
			m = m.toggleAdvancedOptions()
			passToInput = false
			var cmd tea.Cmd
			m.Inputs, cmd = updateFocusStyles(m.Inputs, m.Focus)
			cmds = append(cmds, cmd)

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
			inputChanged = true
			f := &m.Inputs[m.Focus]
			switch f.Type {
			case TypeNumber:
				valStr := m.Inputs[m.Focus].Input.Value()
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					val = 0
				}

				step := 0.1
				switch m.Inputs[m.Focus].Name {
				case "Ratio":
					step = 0.2
				case "Max Slack", "Min Slack":
					step = 0.01
				case "Max Pulley", "Min Pulley":
					step = 1
				}

				if msg.String() == "right" {
					val += step
				} else if msg.String() == "left" {
					val -= step
				}

				f.Input.SetValue(formatFloat(val))
				f.Input.CursorEnd()
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
			inputChanged = true
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
		oldVal := m.Inputs[m.Focus].Input.Value()
		cmd := m.updateInputs(msg)
		cmds = append(cmds, cmd)

		if m.Inputs[m.Focus].Input.Value() != oldVal {
			inputChanged = true
		}
	}

	if inputChanged {
		m.Calculating = true

		if m.CancelCalc != nil {
			m.CancelCalc()
		}

		ctx, cancel := context.WithCancel(context.Background())
		m.CancelCalc = cancel

		c2cVal := m.Inputs[0].Input.Value()
		unitVal := m.Inputs[1].Options[m.Inputs[1].Selected]
		ratioVal := m.Inputs[2].Input.Value()
		useAvailableBelts := m.Inputs[3].Checked
		maxSlack := m.Inputs[4].Input.Value()
		minSlack := m.Inputs[5].Input.Value()
		maxPulley := m.Inputs[6].Input.Value()
		minPulley := m.Inputs[7].Input.Value()

		cmds = append(cmds, updateResults(ctx, c2cVal, ratioVal, unitVal, useAvailableBelts, maxSlack, minSlack, maxPulley, minPulley))
	}

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
	header := title + "\n\n"

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

	if m.Calculating {
		resultsString = m.Spinner.View() + " Calculating..."
	} else if m.ErrorMsg != "" {
		resultsString = errorStyle.Render(m.ErrorMsg)
	} else if len(m.Results) == 0 {
		resultsString = "No valid combination found"
	} else {
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(resultsTableStyle).
			Headers("Pulley 1", "Pulley 2", "Ratio", "Belt Length", "Belt Width", "Slack", "Count")

		for _, res := range m.Results {
			var ratio float64 = float64(res.Pulley2) / float64(res.Pulley1)
			t.Row(
				fmt.Sprintf("%dT", res.Pulley1),
				fmt.Sprintf("%dT", res.Pulley2),
				fmt.Sprintf("%.2f", ratio),
				fmt.Sprintf("%d mm (%dT) ", res.BeltLength, res.BeltLength/5),
				fmt.Sprintf("%.1f mm", res.BeltWidth),
				fmt.Sprintf("%.2f mm", res.Slack),
				fmt.Sprintf("%d", res.Count),
			)
		}

		resultsString = t.Render()
	}

	rightColumnContent := boxTitle + "\n\n" + resultsString

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumnStyle.Render(leftColumn),
		RenderWithMinSize(rightColumnContent, 69, 18),
	)

	footer := helpStyle.Render("\n\nPress `q` to quit, `a` for advanced options, `up/down` to navigate, `left/right` to change values")

	return header + mainContent + footer
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func updateResults(
	ctx context.Context,
	c2cStr string,
	ratioStr string,
	unit string,
	useBelts bool,
	maxSlackString string,
	minSlackString string,
	maxPulleyString string,
	minPulleyString string,
) tea.Cmd {
	return func() tea.Msg {
		results, err := RunCalculator(ctx, c2cStr, ratioStr, unit, useBelts, maxSlackString, minSlackString, maxPulleyString, minPulleyString)

		return CalcResultMessage{
			results: results,
			err:     err,
		}
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

func (m Model) toggleAdvancedOptions() Model {
	for i := range m.Inputs {
		if m.Inputs[i].Advanced {
			m.Inputs[i].Visible = !m.Inputs[i].Visible
		}
	}

	if !m.Inputs[m.Focus].Visible {
		for {
			m.Focus = (m.Focus - 1 + len(m.Inputs)) % len(m.Inputs)
			if m.Inputs[m.Focus].Visible {
				break
			}
		}
	}

	return m
}

func formatFloat(val float64) string {
	s := fmt.Sprintf("%.3f", val)
	s = strings.TrimRight(s, "0")
	s = strings.TrimSuffix(s, ".")
	return s
}
