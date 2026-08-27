package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Belt struct {
	Length float64
	Width  float64
	Count  int
}

type PulleyResult struct {
	Pulley1    int
	Pulley2    int
	BeltLength int
	BeltWidth  float64
	Slack      float64
	Count      int
	Score      float64
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

	var belts []Belt
	if useBelts {
		belts = LoadInventory("belts.csv")
	} else {
		belts = GenerateDefaultBelts()
	}

	minPulley := 8
	maxPulley := 100
	maxSlack := 0.2
	minSlack := -0.5

	slackPenaltyMult := 15
	ratioPenaltyMult := 5

	var pulleyResults []PulleyResult
	for _, belt := range belts {
		for pulley1 := minPulley; pulley1 <= maxPulley; pulley1++ {
			for pulley2 := minPulley; pulley2 <= maxPulley; pulley2++ {
				slack := CalculateSlack(c2c, pulley1, pulley2, belt.Length, 5)

				if slack > maxSlack || slack < minSlack {
					continue
				}
				score := math.Abs(slack) * float64(slackPenaltyMult)

				actualRatio := float64(pulley2) / float64(pulley1)
				ratioError := math.Abs(actualRatio - ratio)
				score += ratioError * float64(ratioPenaltyMult)

				score += calculatePulleySizeScore(pulley1)
				score += calculatePulleySizeScore(pulley2)

				pulleyResults = append(pulleyResults, PulleyResult{
					Pulley1:    pulley1,
					Pulley2:    pulley2,
					BeltLength: int(belt.Length),
					BeltWidth:  belt.Width,
					Slack:      slack,
					Count:      belt.Count,
					Score:      score,
				})
			}
		}
	}

	sort.Slice(pulleyResults, func(i, j int) bool {
		return pulleyResults[i].Score < pulleyResults[j].Score
	})

	if len(pulleyResults) > 10 {
		pulleyResults = pulleyResults[:10]
	}

	return pulleyResults, nil
}

func CalculateSlack(targetC2C float64, pulley1 int, pulley2 int, beltLength float64, pitch float64) float64 {
	beltTeeth := beltLength / pitch

	n1 := float64(pulley1)
	n2 := float64(pulley2)

	// Avoid using math.Pow for squaring numbers as thats significantly slower
	y := beltTeeth - (n1+n2)/2.0
	diff := n2 - n1
	x := (2.0 * (diff * diff)) / (math.Pi * math.Pi)
	radicand := (y * y) - x
	if radicand < 0 {
		return math.MaxFloat64
	}
	actualC2C := (pitch / 4.0) * (y + math.Sqrt(radicand))

	return targetC2C - actualC2C
}

func LoadInventory(filepath string) []Belt {
	file, err := os.Open(filepath)
	if err != nil {
		return GenerateDefaultBelts()
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return GenerateDefaultBelts()
	}

	var belts []Belt
	for i, record := range records {
		if i == 0 && strings.ToLower(record[0]) == "length" {
			continue
		}

		if len(record) < 3 {
			continue
		}

		length, err1 := strconv.ParseFloat(record[0], 64)
		width, err2 := strconv.ParseFloat(record[1], 64)
		count, err3 := strconv.Atoi(record[2])

		if err1 == nil && err2 == nil && err3 == nil && count > 0 {
			belts = append(belts, Belt{Length: length, Width: width, Count: count})
		}
	}

	if len(belts) == 0 {
		return GenerateDefaultBelts()
	}

	return belts
}

func GenerateDefaultBelts() []Belt {
	var belts []Belt
	for teeth := 20; teeth <= 300; teeth++ {
		belts = append(belts, Belt{Length: float64(teeth * 5), Width: 9, Count: 0})
	}
	return belts
}

func calculatePulleySizeScore(pulleySize int) float64 {
	idealMin := 16
	idealMax := 40
	sizePenaltyMult := .5

	var score float64 = 0
	if pulleySize < idealMin {
		score += float64(idealMin-pulleySize) * float64(sizePenaltyMult)
	} else if pulleySize > idealMax {
		score += float64(pulleySize-idealMax) * float64(sizePenaltyMult)
	}

	return score
}
