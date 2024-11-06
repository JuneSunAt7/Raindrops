package main

import (
	"crypto/tls"
	"flag"
	"net"

	"fmt"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"

	client "github.com/JuneSunAt7/Raindrops/1client"
	commands "github.com/JuneSunAt7/Raindrops/cloud_compute"
)

const (
	PORT = "2121"
	HOST = "localhost"
)

func Run() (err error) {

	var connect net.Conn

	boolTSL := flag.Bool("tls", true, "Set tls connection")
	flag.Parse()
	if !*boolTSL {

		connect, err = net.Dial("tcp", HOST+":"+PORT)
		if err != nil {
			pterm.Warning.Println("Не удалось связаться с системой\nИзмените конфигурационный файл или попробуйте снова")
			client.Configure()
			return err
		}

	} else {

		conf := &tls.Config{
			 InsecureSkipVerify: true,
		}

		connect, err = tls.Dial("tcp", HOST+":"+PORT, conf)
		if err != nil {
			return err
		}
	}

	defer connect.Close()

	if err := client.AuthenticateClient(connect); err != nil {

		return err
	}
	
	var options []string

	options = append(options, fmt.Sprintf("Загрузить данные"))
	options = append(options, fmt.Sprintf("Статистика и анализ ранее загруженных данных"))
	options = append(options, fmt.Sprintf("Выход"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.TextStyle.Add(*pterm.NewStyle(pterm.FgMagenta))
	printer.KeyConfirm = keys.Enter

	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		switch selectedOptions {
		case "Загрузить данные":
			commands.UploadBigData(connect)
		case "Статистика и анализ ранее загруженных данных":
			commands.ListStatistics(connect)
		case "Выход":
			client.Exit(connect)
			return
		}
	}

}

func main() {
	Run()
}
