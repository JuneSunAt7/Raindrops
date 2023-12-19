package main

import (
	"fmt"

	"atomicgo.dev/keyboard/keys"
	mgmt "github.com/JuneSunAt7/netMg/server_mgmt"
	"github.com/pterm/pterm"
)

func TuiCommands() {
	var options []string

	options = append(options, fmt.Sprintf("Добавить пользователя"))
	options = append(options, fmt.Sprintf("Удалить пользователя"))
	options = append(options, fmt.Sprintf("Данные о пользователях"))
	options = append(options, fmt.Sprintf("Конфигурация"))
	options = append(options, fmt.Sprintf("Выход"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.TextStyle.Add(*pterm.NewStyle(pterm.FgBlue))
	printer.KeyConfirm = keys.Enter

	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		switch selectedOptions {
		case "Добавить пользователя":
			mgmt.AddUser()
		case "Удалить пользователя":
			mgmt.DelUser()
		case "Данные о пользователях":
			mgmt.UserData()
		case "Конфигурация":
			mgmt.ConfigUser()
		case "Выход":
			return
		}
	}

}

func main() {
	TuiCommands()
}
