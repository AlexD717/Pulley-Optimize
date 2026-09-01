package main

import (
	"context"
	"testing"
)

func TestCalculator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		c2cStr         string
		ratioStr       string
		unit           string
		useBelts       bool
		maxSlack       string
		minSlack       string
		maxPulley      string
		minPulley      string
		expectedResult []PulleyResult
		expectedErr    bool
	}{
		{
			"C2C Not a Number",
			"56Not a Number",
			"3",
			"in",
			false,
			"0.5",
			"-0.2",
			"100",
			"8",
			[]PulleyResult{},
			true,
		},
		{
			"C2C is 0",
			"0",
			"3",
			"in",
			false,
			"0.5",
			"-0.2",
			"100",
			"8",
			[]PulleyResult{},
			true,
		},
		{
			"Ratio Not a Number",
			"5.3",
			"~3",
			"in",
			false,
			"0.5",
			"-0.2",
			"100",
			"8",
			[]PulleyResult{},
			true,
		},
		{
			"Min Pulley Greater than Max Pulley",
			"5.3",
			"~3",
			"in",
			false,
			"0.5",
			"-0.2",
			"100",
			"800",
			[]PulleyResult{},
			true,
		},
		{
			"Min Slack Greater than 0",
			"5.3",
			"~3",
			"in",
			false,
			"0.5",
			"0.2",
			"100",
			"800",
			[]PulleyResult{},
			true,
		},
		{
			"Max Slack Less than 0",
			"5.3",
			"~3",
			"in",
			false,
			"-0.5",
			"-0.2",
			"100",
			"800",
			[]PulleyResult{},
			true,
		},
		{
			"No Result for Small C2C",
			"1",
			"2",
			"mm",
			false,
			"0.5",
			"-0.2",
			"100",
			"20",
			[]PulleyResult{},
			false,
		},
		{
			"Valid Test (mm)",
			"254",
			"3.8",
			"mm",
			false,
			"0.5",
			"-0.2",
			"100",
			"8",
			[]PulleyResult{
				{Pulley1: 10, Pulley2: 38, BeltLength: 630},
				{Pulley1: 9, Pulley2: 37, BeltLength: 625},
			},
			false,
		},
		{
			"Valid Test (in)",
			"10",
			"3.8",
			"in",
			false,
			"0.5",
			"-0.2",
			"100",
			"8",
			[]PulleyResult{
				{Pulley1: 10, Pulley2: 38, BeltLength: 630},
				{Pulley1: 9, Pulley2: 37, BeltLength: 625},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			resultPulleys, resultError := RunCalculator(ctx, tt.c2cStr, tt.ratioStr, tt.unit, tt.useBelts, tt.maxSlack, tt.minSlack, tt.maxPulley, tt.minPulley)

			if (resultError != nil) != tt.expectedErr {
				t.Errorf("RunCalculator error = %v, expectedErr %v", resultError, tt.expectedErr)
				return
			}

			if len(tt.expectedResult) == 0 {
				if len(resultPulleys) != 0 {
					t.Errorf("Expected empty results, got %d results", len(resultPulleys))
				}
			} else {
				for _, expected := range tt.expectedResult {
					found := false
					for _, result := range resultPulleys {
						if result.Pulley1 == expected.Pulley1 && result.Pulley2 == expected.Pulley2 && result.BeltLength == expected.BeltLength {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected result not found in output. Pulley1 %d, Pulley2 %d, BeltLength %d", expected.Pulley1, expected.Pulley2, expected.BeltLength)
					}
				}
			}
		})
	}
}
