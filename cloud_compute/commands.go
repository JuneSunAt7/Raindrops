package cloudcompute
import (
	"net"
	"strings"

	"log"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/pterm/pterm"
)

func UploadBigData(conn net.Conn) {
	fname := ChooseFileData()
	fname = strings.Replace(fname, "\\", "/", -1)

}
func ChooseFileData() string {
	openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
		Title: "Open A File",
		Role:  "OpenFile",
		FileFilters: []cfd.FileFilter{
			{DisplayName: "CSV files", Pattern: "*.csv"},
			{DisplayName: "DB files", Pattern: "*.db"},
			{DisplayName: "XLSX files", Pattern: "*.xlsx"},
		},

	})
	if err != nil {
		log.Fatal(err)
	}

	if err := openDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := openDialog.GetResult()
	if err == cfd.ErrorCancelled {
		pterm.Error.Println("Вы закрыли окно выбора файла")
	} else if err != nil {
		log.Fatal(err)
	}
	return result
}