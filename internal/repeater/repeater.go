package repeater

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AfterNow сравнивает две даты, возвращает true, если дата строго позже now
func AfterNow(date, now time.Time) bool {
	y0, m0, d0 := date.Date()
	y1, m1, d1 := now.Date()

	if y0 != y1 {
		return y0 > y1
	}
	if m0 != m1 {
		return m0 > m1
	}
	return d0 > d1
}

// NextDate возвращает следующую дату повторения задачи, с учетом начальной даты и правил повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect date format: %w", err)
	}
	notEmptyRepeat := strings.TrimSpace(repeat)
	if notEmptyRepeat == "" {
		return "", errors.New("incorrect repeat rule")
	}
	partsRepeat := strings.SplitN(notEmptyRepeat, " ", 2)

	switch partsRepeat[0] {
	case "d":
		return nextDailyDate(now, date, partsRepeat)
	case "w":
		return nextWeeklyDate(now, date, partsRepeat)
	case "m":
		return nextMonthlyDate(now, date, partsRepeat)
	case "y":
		return nextYearlyDate(now, date)
	default:
		return "", fmt.Errorf("incorrect repeat rule: %s", partsRepeat[0])
	}
}

// nextDailyDate возвращает следующую дату для ежедневного правила с заданным интервалом
func nextDailyDate(now time.Time, date time.Time, partsRepeat []string) (string, error) {
	if len(partsRepeat) != 2 {
		return "", errors.New("incorrect daily repeat rule format")
	}
	numberOfDays, err := strconv.Atoi(strings.TrimSpace(partsRepeat[1]))
	if err != nil {
		return "", fmt.Errorf("incorrect number of days in repeat rule: %w", err)
	}
	if numberOfDays < 1 || numberOfDays > 400 {
		return "", fmt.Errorf("number of days must be between 1 and 400, got: %d", numberOfDays)
	}
	date = date.AddDate(0, 0, numberOfDays)
	for !AfterNow(date, now) {
		date = date.AddDate(0, 0, numberOfDays)
	}
	return date.Format("20060102"), nil
}

// nextWeeklyDate возвращает следующую дату для еженедельного правила по дням недели
func nextWeeklyDate(now time.Time, date time.Time, partsRepeat []string) (string, error) {
	if len(partsRepeat) != 2 {
		return "", errors.New("incorrect weekly repeat rule format")
	}
	var weekdays[8] bool
	for _, days := range strings.Split(partsRepeat[1], ",") {
		day, err := strconv.Atoi(strings.TrimSpace(days))
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("incorrect day of week: %d", day)
			}
		weekdays[day] = true
	}
	for {
		dayOfWeek := int(date.Weekday())
			if dayOfWeek == 0 {
				dayOfWeek = 7
			}
		if weekdays[dayOfWeek] && AfterNow(date, now) {
			return date.Format("20060102"), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

// nextMonthlyDate возвращает следующую дату для ежемесячного правила 
func nextMonthlyDate(now time.Time, date time.Time, partsRepeat []string) (string, error) {
	if len(partsRepeat) != 2 {
		return "", errors.New("incorrect monthly repeat rule format")
	}

	mDetails := strings.SplitN(partsRepeat[1], " ", 2)
	mDetailsDay := mDetails[0]
	
	mDetailsMonth := ""
	if len(mDetails) == 2{
		mDetailsMonth = mDetails[1]
	}

	var days [32]bool
	var lastDay, penultimateDay bool

	for _, day := range strings.Split(mDetailsDay, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil || n < -2 || n == 0 || n > 31 {
			return "", fmt.Errorf("incorrect day of month: %s", day)
		}
		switch n{
		case -2:
			penultimateDay = true
		case -1:
			lastDay= true
		default:
			days[n] = true
		}
	}
	var months [13]bool
	if mDetailsMonth == "" {
		for n := 1; n <= 12; n++ {
			months[n] = true
		}
	} else {
		for _, month := range strings.Split(mDetailsMonth, ",") {
			n, err := strconv.Atoi(strings.TrimSpace(month))
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("incorrect month: %s", month)
			}
			months[n] = true
		}
	}

	for {
		year, month, day := date.Date()
		numOfMonth := int(month)

		if mDetailsMonth != "" && !months[numOfMonth] {
			date = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		}

		lastDayofMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

		if (days[day] ||
			(penultimateDay && day == lastDayofMonth - 1) ||
			(lastDay && day == lastDayofMonth)) &&
			AfterNow(date, now) {
			return date.Format("20060102"), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

// nextYearlyDate возвращает следующую дату для ежегодного правила 
func nextYearlyDate(now time.Time, date time.Time) (string, error) {
	year := date.Year() + 1
	month := date.Month()
	day := date.Day()

	if month == time.February && day == 29 && (year % 4 != 0 || (year % 100 == 0 && year % 400 != 0)) {
		date = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
	} else {
		date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	for !AfterNow(date, now) {
		year := date.Year() + 1
		month := date.Month()
		day := date.Day()

		if month == time.February && day == 29 && (year % 4 != 0 || (year % 100 == 0 && year % 400 != 0)) {
			date = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
		} else {
			date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		}
	}
	return date.Format("20060102"), nil
}
