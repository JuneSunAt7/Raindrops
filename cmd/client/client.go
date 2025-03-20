package main

import (
    "crypto/tls"
    "flag"
    "net"
    "fmt"
    "time"

    "atomicgo.dev/keyboard/keys"
    "github.com/pterm/pterm"

    client "github.com/JuneSunAt7/Raindrops/1client"
    plugins "github.com/JuneSunAt7/Raindrops/plugins"
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

    if client.Compare() {
        client.StartServer()

        client.SendNotification("Авторезервирование", "Сейчас день резервирования. Выполнить резервирование?")

        timeout := time.After(30 * time.Second) 
        for {
            select {
            case <-timeout:
                fmt.Println("Время ожидания истекло. Авторезервирование отменено.")
                client.StopServer()
                break
            default:
                if client.UserResponse == "yes" {
                    fmt.Println("Выполняется авторезервирование...")
                    client.AutoSendFiles(connect)
                    client.StopServer()
                    break 
                } else if client.UserResponse == "no" {
                    fmt.Println("Авторезервирование отменено.")
                    client.StopServer()
                    break 
                }
                time.Sleep(500 * time.Millisecond) 
            }
            if client.UserResponse != "" { 
                break
            }
        }
    }

    var options []string

    options = append(options, fmt.Sprintf("Загрузить файл"))
    options = append(options, fmt.Sprintf("Скачать файл"))
    options = append(options, fmt.Sprintf("Список файлов"))
    options = append(options, fmt.Sprintf("Конфигурация"))
    options = append(options, fmt.Sprintf("Сертификаты и пароли"))
    options = append(options, fmt.Sprintf("Авторезервирование"))
    options = append(options, fmt.Sprintf("Плагины"))
    options = append(options, fmt.Sprintf("Моя Kanban-доска"))
    options = append(options, fmt.Sprintf("Выход"))

    printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
    printer.Filter = false
    printer.TextStyle.Add(*pterm.NewStyle(pterm.FgBlue))
    printer.KeyConfirm = keys.Enter

    for {
        selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

        switch selectedOptions {
        case "Загрузить файл":
            client.Upload(connect)
        case "Скачать файл":
            client.Download(connect)
        case "Список файлов":
            client.ListFiles(connect)
        case "Конфигурация":
            client.Configure()
        case "Сертификаты и пароли":
            client.CertPasswd(connect)
        case "Авторезервирование":
            client.Autoreserved()
        case "Плагины":
            plugins.TuiPugins(connect)
        case "Worker":
            client.TuiWorker()
        case "Выход":
            client.Exit(connect)
            return
        }
    }
}

func main() {
    Run()
}