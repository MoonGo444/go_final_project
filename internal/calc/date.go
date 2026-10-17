package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const format = "20060102"

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}

	return false
}

func nextWeekday(now time.Time, date time.Time, numDay int) time.Time {
	var nxtDate time.Time

	if numDay == 7 {
		numDay = 0
	}

	for {
		date = date.AddDate(0, 0, 1)
		if date.Weekday() == time.Weekday(numDay) && date.After(now) {
			nxtDate = date
			break
		}
	}

	return nxtDate
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var interval int

	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(format, dstart)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	rep := strings.Split(repeat, " ")

	if rep[0] != "d" && rep[0] != "y" && rep[0] != "w" && rep[0] != "m" {
		return "", fmt.Errorf("invalid character in repeat: %s", repeat)
	}

	if rep[0] == "d" {
		if len(rep) == 1 || len(rep) > 2 {
			return "", fmt.Errorf("repeat rules violated: %s", repeat)
		}
		interval, err = strconv.Atoi(rep[1])
		if err != nil {
			return "", fmt.Errorf("repeat conv: %w", err)
		}
		if interval > 400 || interval < 1 {
			return "", errors.New("repeat out of range")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	}

	for rep[0] == "y" {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	if rep[0] == "w" {
		if len(rep) == 1 || len(rep) > 2 {
			return "", fmt.Errorf("repeat rules violated: %s", repeat)
		}

		repNum := strings.Split(rep[1], ",")
		slTimeDate := make([]time.Time, 0)

		for _, numDay := range repNum {
			numDay, err := strconv.Atoi(numDay)
			if err != nil {
				return "", fmt.Errorf("repeat conv: %w", err)
			}
			if numDay < 1 || numDay > 7 {
				return "", fmt.Errorf("out of range repeat, w: %d", numDay)
			}

			slTimeDate = append(slTimeDate, nextWeekday(now, date, numDay))
		}

		var minTime time.Time

		for i, timeDay := range slTimeDate {
			if i == 0 {
				minTime = timeDay
				continue
			}
			if minTime.After(timeDay) {
				minTime = timeDay
			}
		}

		date = minTime
	}

	if rep[0] == "m" {
		if len(rep) < 2 || len(rep) > 3 {
			return "", fmt.Errorf("repeat rules violated: %s", repeat)
		}

		var day [32]bool
		var month [13]bool

		if len(rep) == 2 {
			for i := range month {
				if i == 0 {
					continue
				}
				month[i] = true
			}

			slDay := strings.Split(rep[1], ",")
			slDate := make([]time.Time, 0)
			var minDate time.Time

			for _, d := range slDay {
				numDay, err := strconv.Atoi(d)
				if err != nil {
					return "", fmt.Errorf("repeat conv: %w", err)
				}
				if numDay < -2 || numDay > 31 || numDay == 0 {
					return "", fmt.Errorf("out of range repeat, w: %d", numDay)
				}

			loop:
				for i := range day {
					if numDay == -1 || numDay == -2 {
						newDate := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
						newDate = newDate.AddDate(0, 0, numDay)

						if newDate.Equal(date) || newDate.Before(date) {
							newDate = time.Date(newDate.Year(), newDate.Month()+2, 1, 0, 0, 0, 0, time.UTC)
							newDate = newDate.AddDate(0, 0, numDay)
						}

						if newDate.Before(now) {
							for {
								newDate = newDate.AddDate(0, 1, 0)
								if afterNow(newDate, now) {
									break
								}
							}
						}

						day[newDate.Day()] = true

						slDate = append(slDate, newDate)

						break loop
					}
					if i == numDay {
						day[i] = true
						var newDate time.Time = date
						for {
							newDate = newDate.AddDate(0, 0, 1)

							if afterNow(newDate, now) {
								if newDate.Day() == i {
									break
								}
							}
						}
						slDate = append(slDate, newDate)
					}
				}
			}

			for i, date := range slDate {
				if i == 0 {
					minDate = date
					continue
				}
				if date.Before(minDate) {
					minDate = date
				}
			}

			if len(slDate) > 0 {
				date = minDate
			}
		}

		if len(rep) == 3 {
			slDay := strings.Split(rep[1], ",")
			slMonth := strings.Split(rep[2], ",")
			slDate := make([]time.Time, 0)

			for _, m := range slMonth {
				numMonth, err := strconv.Atoi(m)
				if err != nil {
					return "", fmt.Errorf("repeat conv: %w", err)
				}
				if numMonth < 1 || numMonth > 12 {
					return "", fmt.Errorf("out of range repeat, w: %d", numMonth)
				}

				for i := range month {
					if i == numMonth {
						month[i] = true
					}
				}
			}

			for _, d := range slDay {
				numDay, err := strconv.Atoi(d)
				if err != nil {
					return "", fmt.Errorf("repeat conv: %w", err)
				}
				if numDay < -2 || numDay > 31 || numDay == 0 {
					return "", fmt.Errorf("out of range repeat, w: %d", numDay)
				}

				if numDay == -1 || numDay == -2 {
					for i, m := range month {
						if m == true {
							newDate := time.Date(date.Year(), time.Month(i)+1, 1, 0, 0, 0, 0, time.UTC)
							newDate = newDate.AddDate(0, 0, numDay)

							if newDate.Equal(date) {
								continue
							}

							if newDate.Before(date) {
								newDate = time.Date(newDate.Year()+1, newDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
								newDate = newDate.AddDate(0, 0, numDay)
								slDate = append(slDate, newDate)
								continue
							}

							slDate = append(slDate, newDate)
						}
					}
				}

				for i := range day {
					if i == numDay {
						day[i] = true
					}
				}
			}

			for i, d := range day {
				if d == true {
					for y, m := range month {
						if m == true {
							monthDays := time.Date(date.Year(), time.Month(y)+1, 0, 0, 0, 0, 0, time.UTC)
							if i > monthDays.Day() && len(slMonth) == 1 && len(slDay) == 1 {
								return "", errors.New("month day out of range")
							} else if i > monthDays.Day() {
								continue
							}

							newDate := time.Date(date.Year(), time.Month(y), i, 0, 0, 0, 0, time.UTC)

							if date.Before(newDate) && newDate.After(now) {
								slDate = append(slDate, newDate)
								continue
							}

							if newDate.Before(date) {
								for {
									newDate = newDate.AddDate(1, 0, 0)
									if newDate.After(now) {
										break
									}
								}
							}

							if newDate.After(date) && newDate.Before(now) {
								for {
									newDate = newDate.AddDate(1, 0, 0)
									if newDate.After(now) {
										break
									}
								}
							}

							if newDate.Equal(date) {
								newDate = newDate.AddDate(1, 0, 0)
							}

							slDate = append(slDate, newDate)
						}
					}
				}
			}

			var minDate = now

			for i, date := range slDate {
				if i == 0 {
					minDate = date
					continue
				}

				if date.Before(minDate) {
					minDate = date
					continue
				}
			}

			date = minDate
		}

	}

	return date.Format(format), nil
}
