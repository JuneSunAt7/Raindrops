package client
import (
    "net/http"
    "github.com/pterm/pterm"
)

func NewWorker(){
	fs := http.FileServer(http.Dir("webworker"))
    http.Handle("/", fs)

    http.ListenAndServe(":80", nil)
    pterm.Success.Println("Запущен клиент на http://localhost:80/")

}
