package client

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"os"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

var PORT string
var HOST string

type ConfigData struct {
	IP   string `json:"ip"`
	PORT string `json:"port"`
}

type NewData struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

func createConfig(ip1, port1 string) {
	user := &NewData{Ip: ip1,
		Port: port1}

	b, err := json.Marshal(user)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println("errr")
	}

	err = ioutil.WriteFile("confRD.conf", b, 0644)

}

func Configure() {
	var options []string

	options = append(options, fmt.Sprintf("Изменить конфигурацию"))
	options = append(options, fmt.Sprintf("Удалить конфигурацию"))
	options = append(options, fmt.Sprintf("Сконфигурировать"))
	options = append(options, fmt.Sprintf("Назад"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter

	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		switch selectedOptions {
		case "Изменить конфигурацию":
			ip, _ := pterm.DefaultInteractiveTextInput.Show("IP")
			port, _ := pterm.DefaultInteractiveTextInput.Show("PORT")
			createConfig(ip, port)
			pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgGreen)).Println("Конфигурация изменена и сохранена")
		case "Удалить конфигурацию":
			createConfig("", "")
			pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgGreen)).Println("Конфигурация очищена")
		case "Сконфигурировать":
			ReadHOST()
			ReadPORT()
		case "Назад":
			return
		}
	}
}

func ReadHOST() string {
	file, _ := os.Open("confRD.conf")
	decoder := json.NewDecoder(file)
	config := new(ConfigData)
	err := decoder.Decode(&config)
	if err != nil {
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgRed)).Println("Отсутствует файл конфигурации")
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgBlue)).Println("Применяется стандартная конфигурация: HOST - localhost, PORT - 2121...")
		HOST = "localhost"
	}
	HOST = config.IP
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgGreen)).Printfln("Успешно сконфигурировано")
	return HOST
}

func ReadPORT() string {
	file, _ := os.Open("confRD.conf")
	decoder := json.NewDecoder(file)
	config := new(ConfigData)
	err := decoder.Decode(&config)
	if err != nil {
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgRed)).Println("Отсутствует файл конфигурации")
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgBlue)).Println("Применяется стандартная конфигурация: HOST - localhost, PORT - 2121...")
		PORT = "2121"
	}
	PORT = config.PORT
	return PORT
}
