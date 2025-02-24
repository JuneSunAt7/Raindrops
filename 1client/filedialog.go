package client

import (
	"log"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/pterm/pterm"
)

func ChooseFile() string {
    openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
        Title: "Open A File",
        Role:  "OpenFile",
    })
    if err != nil {
        log.Fatalf("Ошибка при создании диалогового окна: %v", err)
    }

    if err := openDialog.Show(); err != nil {
        log.Fatalf("Ошибка при показе диалогового окна: %v", err)
    }

    result, err := openDialog.GetResult()
    if err == cfd.ErrorCancelled {
        pterm.Error.Println("Вы закрыли окно выбора файла")
        return ""
    } else if err != nil {
        log.Fatalf("Ошибка при получении результата: %v", err)
    }

    return result
}
func ChooseDir() string {
	openDialog, err := cfd.NewSelectFolderDialog(cfd.DialogConfig{
		Title: "Выбор папки",
		Role:  "OpenFolder",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := openDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := openDialog.GetResult()
	if err == cfd.ErrorCancelled {
		pterm.Error.Println("Вы закрыли окно выбора папки")
	} else if err != nil {
		log.Fatal(err)
	}
	return result
}
