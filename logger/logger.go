package logger

import (
	"log"
	"os"

	"github.com/pterm/pterm"
)

var (
	infoLogger2File *log.Logger
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var file, err1 = os.Create(pwd + "/info.log")
	if err1 != nil {
		panic(err1)
	}

	infoLogger2File = log.New(file, "INFO ", log.LstdFlags) 
}

func Println(v ...interface{}) {
	pterm.Info.Println(v...)
	infoLogger2File.Println(v...)
}
