package client

import (
	"strings"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/sys/windows/registry"
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

	key := registry.CURRENT_USER
	subKey := "Software\\Raindrops"
	valueName := "days"

		val, err := ReadRegistryValue(key, subKey, valueName)
		if err != nil {
			pterm.Error.Println("Error reading registry value:", err)
			return false
		}
	lines := strings.Fields(val)

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
