package editor
import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
)

func Compile(){
	// Компилируем плагин
	err := buildPlugin()
	if err != nil {
		fmt.Println("Ошибка компиляции плагина:", err)
		return
	}

	// Загружаем плагин
	p, err := plugin.Open("plugin.so")
	if err != nil {
		fmt.Println("Ошибка загрузки плагина:", err)
		return
	}

	// Вызываем экспортированную функцию "Hello" из плагина
	helloFunc, err := p.Lookup("Hello")
	if err != nil {
		fmt.Println("Ошибка поиска функции:", err)
		return
	}

	helloFunc.(func())()
}

func buildPlugin()error {
	// Создаем файл с исходным кодом плагина
	file, err := os.Create("plugin.go")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(`
package main

import "fmt"

func Hello() {
	fmt.Println("Hello, Plugin!")
}
`)
	if err != nil {
		return err
	}

	// Компилируем плагин
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "plugin.so", "plugin.go")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Удаляем исходный код плагина
	err = os.Remove("plugin.go")
	if err != nil {
		return err
	}

	return nil
}