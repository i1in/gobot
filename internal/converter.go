package internal

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type TimeConverter struct {
	time  string
	isInt bool
}

func (t *TimeConverter) plural(num int, titles ...string) string {
	if num%100 > 4 && num%100 < 20 {
		return titles[2]
	} else {
		cases := []int{2, 0, 1, 1, 1, 2}
		key := int(math.Min(float64(num%10), 5))
		return titles[cases[key]]
	}
}

func (t *TimeConverter) format(num int, titles ...string) string {
	return fmt.Sprintf("%d %s", num, t.plural(num, titles...))
}

func (t *TimeConverter) timeToDate() string {

	time, err := strconv.Atoi(t.time)
	if err != nil {
		log.Fatal("%w", err)
	}

	days := time / (3600 * 24)
	hours := time / 3600
	minutes := time / 60 % 60

	if days > 0 {
		return t.format(days, "день", "дня", "дней")
	}

	if hours > 0 {
		return t.format(hours, "час", "часа", "часов")
	}

	if minutes > 0 {
		return t.format(minutes, "минута", "минуты", "минут")
	} else {
		return t.format(time, "секунда", "секунды", "секунд")
	}
}

func (t *TimeConverter) convertTime() *TimeConverter {

	var unit = t.time[len(t.time)-1:]

	value, err := strconv.Atoi(strings.TrimSuffix(t.time, unit))
	if err != nil {
		return &TimeConverter{
			time:  t.time,
			isInt: false,
		}
	}

	switch unit {
	case "s":
		return &TimeConverter{
			time:  strconv.Itoa(value),
			isInt: true,
		}
	case "m":
		return &TimeConverter{
			time:  strconv.Itoa(value * 60),
			isInt: true,
		}
	case "h":
		return &TimeConverter{
			time:  strconv.Itoa(value * 60 * 60),
			isInt: true,
		}
	case "d":
		return &TimeConverter{
			time:  strconv.Itoa(value * 24 * 60 * 60),
			isInt: true,
		}
	default:
		return &TimeConverter{
			time:  t.time,
			isInt: false,
		}
	}
}

func convert(time string) *TimeConverter {
	return &TimeConverter{
		time: time,
	}
}
