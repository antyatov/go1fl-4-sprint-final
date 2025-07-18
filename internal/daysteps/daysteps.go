package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) < 2 {
		return 0, 0, fmt.Errorf("reсeived slice has fewer than 2 items")
	}

	stepsCount, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parce steps count: %w", err)
	}

	walkDuration, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse walk duration: %w", err)
	}

	if stepsCount == 0 || walkDuration == 0 {
		return 0, 0, fmt.Errorf("steps count (%d) or wall duration (%d) cannot will be zero: %w", stepsCount, walkDuration, err)
	}

	return stepsCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Printf("failed receive steps and walk duration: %s\n", err)
		return ""
	}

	if stepsCount == 0 {
		log.Printf("steps count equal zero: %d\n", stepsCount)
		return ""
	}

	walkDistance := (float64(stepsCount) * stepLength) / mInKm
	caloriesSpent, err := spentcalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)

	if err != nil {
		log.Printf("failed receive spent calories (%f): %s\n", caloriesSpent, err)
	}

	dayInfoMessageTemplate := `Количество шагов: %d.
		Дистанция составила %f км.
		Вы сожгли %f ккал.`

	return fmt.Sprintf(dayInfoMessageTemplate, stepsCount, walkDistance, caloriesSpent)
}
