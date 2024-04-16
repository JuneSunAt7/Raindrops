package client
import (
    "net/http"
)

func NewWorker(){
	fs := http.FileServer(http.Dir("webworker"))
    http.Handle("/", fs)

    http.ListenAndServe(":8080", nil)
	

}

func ChooseWorker(){

}

func SettingWorker(){

}