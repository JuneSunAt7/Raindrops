package plugins

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"github.com/pterm/pterm"
	"net"
	"fmt"
	"io/ioutil"
	"atomicgo.dev/keyboard/keys"
	"io"
	"bytes"
	"strconv"
	"time"
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
	pterm.BgLightYellow.Printfln("Успешно добавлены плагины:\n %s", pluginsPath)
}

func SearchPluginsInServer(conn net.Conn){
	conn.Write([]byte("pluginshop\n"))
	buffer := make([]byte, 4096)
	n, _ := conn.Read(buffer)

	pterm.BgLightMagenta.Println("Магазин плагинов")
	plugins := strings.Split(string(buffer[:n]), "\n")

	var options []string

	for i := 0; i < len(plugins); i++ {
		options = append(options,fmt.Sprintf(plugins[i]))
	}
	options = append(options, fmt.Sprintf("Назад"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		if selectedOptions == "Назад"{
			return
		}
		getPlugin(conn, selectedOptions)
		
	}
}
func getPlugin(conn net.Conn, plugin string){
	conn.Write([]byte(fmt.Sprintf("getplugin %s\n",plugin)))
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	comStr := strings.Trim(string(buffer[:n]), "\n")
	commandArr := strings.Fields(comStr)

	fileSize, err := strconv.ParseInt(commandArr[2], 10, 64)
	if err != nil || fileSize == -1 {
		pterm.Error.Println("Ошибка при скачивании плагина")
		conn.Write([]byte("file size error"))
		return
	}
	conn.Write([]byte("200 Start download!"))
	var arrDec []byte

	outputFile, err := os.Create("plugins/"+ plugin)
	if err != nil {
		pterm.Error.Println("Ошибка при скачивании плагина")
		
	}
	io.Copy(outputFile, bytes.NewReader(arrDec))
	defer outputFile.Close()
	p, _ := pterm.DefaultProgressbar.WithTotal(3).WithTitle("...Скачивание плагина...").Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle("Загрузка из магазина") // ProgressBar - downloader
		p.Increment()
		time.Sleep(time.Millisecond * 350)
		
	}
	pterm.Success.Println("Успешная загрузка плагина!")
	var plugPath []string
	plugPath = append(plugPath,  plugin)
	AddPlugin(plugPath)
}
func RunPlugin(){
	filePath := "plugins/plugins.ini"

    
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
		pterm.Error.Println("Ошибка запуска, не обнаружено ни одного плагина.\nСкачайте их из магазина или найдите на компьютере")
        return
    }
    lines := strings.Split(string(content), "\n")


    var arr []string
    for _, line := range lines {
        arr = append(arr, line)
    }
	InitPlugins(arr)
}