package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		activity := ""
		if len(parts) > 1 {
			activity = parts[1]
		}
		return 0, activity, 0, errors.New("invalid training data format")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, parts[1], 0, err
	}
	if steps <= 0 {
		return 0, parts[1], 0, errors.New("steps must be positive")
	}

	dur, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, parts[1], 0, err
	}
	if dur <= 0 {
		return 0, parts[1], 0, errors.New("duration must be positive")
	}

	return steps, parts[1], dur, nil
}

func distance(steps int, height float64) float64 {
	lenStep := height * stepLengthCoefficient
	distance := (float64(steps) * lenStep) / mInKm

	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

// TrainingInfo parses the input data string, calculates distance, average speed,
// and calories burned for the specified activity, and returns a formatted result string.
//
// Parameters:
//
//	data string      — input string in the format "3456,Ходьба,3h00m",
//	                   containing the number of steps, activity type, and duration
//	weight float64   — user's weight in kilograms
//	height float64   — user's height in meters
//
// Returns:
//
//	string — formatted string with workout information, e.g.:
//
//	   Тип тренировки: Бег
//	   Длительность: 0.75 ч.
//	   Дистанция: 10.00 км.
//	   Скорость: 13.34 км/ч
//	   Сожгли калорий: 18621.75
//
//	error  — non-nil if input parsing fails or if an unknown activity type is provided
//
// Behavior:
//   - Calls parseTraining() to extract steps, activity, and duration.
//   - Logs parsing errors using log.Println(err) and returns "Error" with the error.
//   - Calculates distance, average speed, and calories based on activity type
//     (Walking or Running).
//   - Returns an error if the activity type is unknown.
func TrainingInfo(data string, weight, height float64) (string, error) {
	var calories float64
	var errCal error

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := distance / duration.Hours()

	switch activity {
	case "Ходьба":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if errCal != nil {
		log.Println(errCal)
		return "", errCal
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distance,
		speed,
		calories,
	), nil
}

// RunningSpentCalories calculates calories burned during running.
//
// It takes the following parameters:
//
//	steps int       — number of steps taken
//	weight float64  — user's weight in kilograms
//	height float64  — user's height in meters
//	duration time.Duration — running duration
//
// It returns two values:
//
//	float64 — calories burned during the run
//	error   — non-nil if input parameters are invalid (e.g., non-positive steps, weight, height, or duration)
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}

	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}

	if height <= 0 {
		return 0, errors.New("height must be positive")
	}

	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil
}

// WalkingSpentCalories calculates the number of calories burned during walking.
//
// Parameters:
//
//	steps int       — number of steps taken
//	weight float64  — user's weight in kilograms
//	height float64  — user's height in meters
//	duration time.Duration — duration of the walk
//
// Returns:
//
//	float64 — calories burned during the walk
//	error   — non-nil if input parameters are invalid (e.g., non-positive steps, weight, height, or duration)
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}

	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}

	if height <= 0 {
		return 0, errors.New("height must be positive")
	}

	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	return ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}
