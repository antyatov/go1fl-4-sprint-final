package spentcalories

import (
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
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) < 3 {
		return 0, "", 0, fmt.Errorf("reсeived slice has fewer than 2 items")
	}

	stepsCount, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("cannot parce steps count: %w", err)
	}

	activityKind := dataSlice[1]

	activityDuration, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("cannot parce duration of activity: %w", err)
	}

	return stepsCount, activityKind, activityDuration, nil
}

func distance(steps int, height float64) float64 {
	distanceMeters := (height * stepLengthCoefficient) * float64(steps)
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	activityMessage := `Тип тренировки: %s
		Длительность: %s ч.
		Дистанция: %f км.
		Скорость: %f км/ч
		Сожгли калорий: %f`

	switch strings.ToLower(activity) {
	case "ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("failed calc calories: %w", err)
		}

		return fmt.Sprintf(activityMessage, activity, duration, dist, speed, calories), nil

	case "бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("failed calc calories: %w", err)
		}

		return fmt.Sprintf(activityMessage, activity, duration, dist, speed, calories), nil

	default:
		return "", fmt.Errorf("unknown activity kind\n(неизвестный тип тренировки):\n%s", activity)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps value :%d", steps)
	}

	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("invalid weight(%f) or height (%f) value", weight, height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid diration time: %s", duration)
	}

	runningSpeed := meanSpeed(steps, height, duration)
	calories := (weight * runningSpeed * duration.Minutes())

	return calories / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps value :%d", steps)
	}

	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("invalid weight(%f) or height (%f) value", weight, height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("invalid diration time: %s", duration)
	}

	walkingSpeed := meanSpeed(steps, height, duration)
	calories := (weight * walkingSpeed * duration.Minutes())

	return (calories / minInH) * walkingCaloriesCoefficient, nil
}
