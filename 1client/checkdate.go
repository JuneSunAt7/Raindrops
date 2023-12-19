package client

import (
	"bufio"
	"os"
	"time"

	"github.com/pterm/pterm"
)

var days = map[string]string{
	"Понедельник": "Monday",
	"Вторник":     "Tuesday",
	"Среда":       "Wensday",
	"Четверг":     "Thursday",
	"Пятница":     "Friday",
	"Суббота":     "Saturday",
	"Воскресенье": "Sunday",
}

func todayData() string {
	now := time.Now()
	weekday := now.Weekday().String()
	return weekday
}

func Compare() bool {
	var reserve bool

	file, err := os.Open(ROOT + "/" + "localSettings" + "/" + "settings.ini")
	if err != nil {
		pterm.FgRed.Println("Файл настроек не найден")
		return false
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < 7; i++ {
		for j := 0; j < len(lines); j++ {
			if days[lines[j]] == todayData() {
				return true
			} else {
				reserve = false
			}
		}
	}

	return reserve
}
