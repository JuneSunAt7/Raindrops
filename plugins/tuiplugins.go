package plugins
import(
	
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"fmt"
	masterdir "github.com/JuneSunAt7/Raindrops/1client"
)
func TuiPugins(){
	var options []string

	options = append(options, fmt.Sprintf("Доступные плагины"))
	options = append(options, fmt.Sprintf("Запустить плагин"))
	options = append(options, fmt.Sprintf("Магазин плагинов"))
	options = append(options, fmt.Sprintf("Найти на этом компьютере плагины"))
	options = append(options, fmt.Sprintf("Назад"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		switch selectedOptions {
		case "Доступные плагины":

	
		case "Запустить плагин":
	
		case "Магазин плагинов":
	
		case "Найти на этом компьютере плагины":
			currentDir := masterdir.ChooseDir()
			CheckPlugins(currentDir, ".so")
		case "Назад":
			return
		}
	}
}