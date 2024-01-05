package plugins

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"github.com/pterm/pterm"
	"net"

)
// Interface for plugins
type Plugin interface {
	Run()
}
func ReadLoacalPlugins(){
	CheckPlugins("plugins", ".so")
}
func CheckPlugins(root, ext string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			pterm.Error.Println("Ошибка поиска плагинов")
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0{
		pterm.Warning.Printfln("В папке %s нет подходящих плагинов", root)
	}else{
		AddPlugin(files)
	}
	return files, nil
}

func InitPlugins(plugins []string) {
	for i :=0; i < len(plugins); i++{
		path :=plugins[i]

	p, err := plugin.Open(path)
	if err != nil {
		pterm.Error.Println("Ошибка загрузки плагина")
		return
	}
	// Get link to Run()
	runSym, err := p.Lookup("Run")
	if err != nil {
		pterm.Error.Println("Ошибка поиска функции `Run` в плагине")
		return
	}
	// Check link to interface
	var pluginFunc Plugin
	pluginFunc, ok := runSym.(Plugin)
	if !ok {
		pterm.Error.Println("Функция `Run` не реализует требуемый интерфейс")
		return
	}
	pluginFunc.Run()
	}
}
func AddPlugin(pluginsPath[] string) {
	outputFile, err := os.Create("plugins/plugins.ini")
	if err != nil {
		pterm.FgRed.Printfln("Ошибка добавления плагина")
		return
	}
	defer outputFile.Close()
	outputFile.WriteString(strings.Join(pluginsPath, "\n"))
	pterm.Success.Printfln("Успешно добавлены плагины\n %s", pluginsPath)
}
func SearchPluginsInServer(conn net.Conn){
	conn.Write([]byte("plugin shop"))
	
}