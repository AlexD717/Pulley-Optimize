package main

import (
	"fmt"
	"strconv"
)

type PulleyResult struct {
	Pulley1     int
	Pulley2     int
	BeltLength  int
	BeltWidth   float64
	Slack       float64
	IsAvailable bool
	Score       float64
}

func RunCalculator(c2cStr string, ratioStr string, unit string, useBelts bool) ([]PulleyResult, error) {
	c2c, err := strconv.ParseFloat(c2cStr, 64)
	if err != nil || c2c <= 0 {
		return []PulleyResult{}, fmt.Errorf("Invalid C2C Distance: must be a number greater than 0")
	}

	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil || ratio <= 0 {
		return []PulleyResult{}, fmt.Errorf("Invalid Ratio: must be a number greater than 0")
	}

	if unit == "in" {
		c2c = c2c * 25.4
		unit = "mm"
	}

	dummyResult := PulleyResult{
		Pulley1:     24,
		Pulley2:     48,
		BeltLength:  250,
		BeltWidth:   10,
		Slack:       0.05,
		IsAvailable: false,
		Score:       0,
	}

	dummyResult2 := PulleyResult{
		Pulley1:     36,
		Pulley2:     48,
		BeltLength:  300,
		BeltWidth:   15,
		Slack:       1,
		IsAvailable: true,
		Score:       0,
	}

	return []PulleyResult{dummyResult, dummyResult2}, nil
}
