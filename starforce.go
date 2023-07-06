package main

import (
	"fmt"
	"math"
	"math/rand"
	"gonum.org/v1/gonum/stat"
	"sort"
)

func main() {
	fmt.Println(mesosPerAttempt(17, 200))
	fmt.Println(getPercentiles(17, 20, 200))
}

func mesosPerAttempt(currentStar, equipLevel int) int {
	var multiplier float64
	var exponent float64
	switch {
	case currentStar < 10:
		exponent = 1
	default:
		exponent = 2.7
	}
	switch {
	case currentStar < 10:
		multiplier = 2500
	case currentStar == 10:
		multiplier = 40000
	case currentStar == 11:
		multiplier = 22000
	case currentStar == 12:
		multiplier = 15000
	case currentStar == 13:
		multiplier = 11000
	case currentStar == 14:
		multiplier = 7500
	default:
		multiplier = 20000
	}
	return int(100 * math.Round(math.Pow(float64(equipLevel), 3)*math.Pow((float64(currentStar)+1), exponent) / multiplier +10))
}

func getSaviorRates() map[int][]float64 {
	// Success, Fail(keep), Fail(decrease), Fail(Destroy)
	rates := map[int][]float64{
		0: []float64{0.95, 0.05, 0, 0},
    1: []float64{0.9, 0.1, 0, 0},
    2: []float64{0.85, 0.15, 0, 0},
    3: []float64{0.85, 0.15, 0, 0},
    4: []float64{0.80, 0.2, 0, 0},
    5: []float64{0.75, 0.25, 0, 0},
    6: []float64{0.7, 0.3, 0, 0},
    7: []float64{0.65, 0.35, 0, 0},
		8: []float64{0.6, 0.4, 0, 0},
    9: []float64{0.55, 0.45, 0, 0},
    10: []float64{0.5, 0.5, 0, 0},
		11: []float64{0.45, 0.55, 0.0, 0.0},
    12: []float64{0.4, 0.6, 0.0, 0.0},
    13: []float64{0.35, 0.65, 0.0, 0.0},
    14: []float64{0.3, 0.7, 0.0, 0.0},
    15: []float64{0.3, 0.679, 0, 0.021},
    16: []float64{0.3, 0.0, 0.679, 0.021},
    17: []float64{0.3, 0.0, 0.679, 0.021},
    18: []float64{0.3, 0.0, 0.672, 0.028},
    19: []float64{0.3, 0.0, 0.672, 0.028},
    20: []float64{0.3, 0.63, 0, 0.07},
    21: []float64{0.3, 0, 0.63, 0.07},
    22: []float64{0.03, 0.0, 0.776, 0.194},
    23: []float64{0.02, 0.0, 0.686, 0.294},
    24: []float64{0.01, 0.0, 0.594, 0.396},
	}
	return rates
}

func getStarforceAttemptOutcome(currentStars int, consecutiveDecreases int) string{
	if consecutiveDecreases == 2{
		return "success"
	}
	rates := getSaviorRates()[currentStars]
	attempt := rand.Float64()
	if attempt <= rates[0] {
		return "success"
	}
	if attempt <= rates[0] + rates[1] {
		return "keep"
	}
	if attempt <= rates[0] + rates[1] + rates[2] {
		return "decrease"
	}
	return "boom"
}

func simulateOneAttempt(currentStars, targetStars, equipLevel int) (int, int){
	numBooms := 0
	consecutiveDecreases := 0
	mesosUsed := 0
	for currentStars < targetStars {
		mesosUsed += mesosPerAttempt(currentStars, equipLevel)
		starforceOutcome := getStarforceAttemptOutcome(currentStars, consecutiveDecreases)
		switch starforceOutcome {
		case "success":
			currentStars += 1
			consecutiveDecreases = 0
		case "keep":
			consecutiveDecreases = 0
		case "decrease":
			currentStars -= 1
			consecutiveDecreases += 1
		case "boom":
			currentStars = 12
			consecutiveDecreases = 0 
			numBooms += 1
		}
	}
	return mesosUsed, numBooms
}

func getPercentiles(currentStars, targetStars, equipLevel int) ([]int, []int) {
	numRuns := 10000
	mesosUsedArray := make([]float64, numRuns)
	numBoonsArray := make([]float64, numRuns)
	for i:=0; i < 10000; i++ {
		mesosUsed, numBoons := simulateOneAttempt(currentStars, targetStars, equipLevel)
		mesosUsedArray[i] = float64(mesosUsed)
		numBoonsArray[i] = float64(numBoons)
	}
	sort.Float64s(mesosUsedArray)
	sort.Float64s(numBoonsArray)
	percentiles := []float64{0.50, 0.75, 0.90}
	mesosUsedPercentiles := make([]int, len(percentiles))
	numBoonsPercentiles := make([]int, len(percentiles))
	for i, p := range percentiles {
		mesosUsedPercentiles[i] = int(stat.Quantile(p, stat.Empirical, mesosUsedArray, nil))
		numBoonsPercentiles[i] = int(stat.Quantile(p, stat.Empirical, numBoonsArray, nil))
	}
	return mesosUsedPercentiles, numBoonsPercentiles
}
