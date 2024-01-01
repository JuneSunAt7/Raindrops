package plugins

import(
	"os"
	"path/filepath"
	"github.com/pterm/pterm"

)
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
	return files, nil
}

func InitPlugins(plugins []string){
	
}
func AddPlugin(path string){

}