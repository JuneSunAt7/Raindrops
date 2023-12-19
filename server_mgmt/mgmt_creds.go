package servermgmt

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"atomicgo.dev/keyboard/keys"

	"github.com/pterm/pterm"
)

type MyStruct struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}
type CredArr []MyStruct

func (p *CredArr) FromJSON(r io.Reader) error {
	en := json.NewDecoder(r)
	return en.Decode(p)
}
func GetCred() (*CredArr, error) {
	f, err := os.Open("user_creds.db")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var creds CredArr
	err = creds.FromJSON(f)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}
func AddUser() {
	uname, _ := pterm.DefaultInteractiveTextInput.Show("Имя")
	passwd, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Пароль")

	filename := "user_creds.db"
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	data := []MyStruct{}

	json.Unmarshal(file, &data)

	newStruct := &MyStruct{
		UserName: uname,
		PassWord: passwd,
	}

	data = append(data, *newStruct)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile(filename, dataBytes, fs.FileMode(0644))
	if err != nil {
		fmt.Println(err)
	}

}
func DelUser() {
	creds, err := GetCred()
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	var options []string
	for _, cred := range *creds {
		options = append(options, fmt.Sprintf(cred.UserName))

	}
	options = append(options, fmt.Sprintf("Назад"))
	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.TextStyle.Add(*pterm.NewStyle(pterm.FgBlue))
	printer.KeyConfirm = keys.Enter

	selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	if selectedOptions == "Назад" {
		return
	} else {
		pterm.Info.Printfln("Выбранный пользователь: %s", pterm.Green(selectedOptions))
		delUserinCreds(selectedOptions)
	}

}
func delUserinCreds(key string) {
	file, err := os.Open("user_creds.db")
	if err != nil {
		pterm.Error.Println("Ошибка при открытии файла")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		pterm.Error.Println("Ошибка при чтении файла")
		return
	}

	var records []MyStruct
	err = json.Unmarshal(data, &records)
	if err != nil {
		pterm.Error.Println("Ошибка при парсинге")
		return
	}

	records = removeRecord(records, key)

	newData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		pterm.Error.Println("Ошибка при кодировании данных")
		return
	}

	err = ioutil.WriteFile("user_creds.db", newData, 0644)
	if err != nil {
		pterm.Error.Println("Ошибка при записи в файл")
		return
	}
	pterm.Success.Println("Пользователь ", key, " успешно удален")
}
func removeRecord(records []MyStruct, name string) []MyStruct {
	var updatedRecords []MyStruct
	for _, record := range records {
		if record.UserName != name {
			updatedRecords = append(updatedRecords, record)
		}
	}
	return updatedRecords
}
func UserData() {
	dirPath := "filestore/"
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirSize, err := getDirSize(path)
			if err != nil {
				return err
			}
			bars := []pterm.Bar{
				{Label: path, Value: int(dirSize)},
			}
			pterm.DefaultBarChart.WithBars(bars).WithHorizontal().WithShowValue().Render()
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

}
func getDirSize(dirPath string) (int64, error) {
	var size int64

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if !file.IsDir() {
			size += file.Size()
		} else {
			subDirSize, err := getDirSize(filepath.Join(dirPath, file.Name()))
			if err != nil {
				return 0, err
			}
			size += subDirSize
		}
	}

	return size, nil
}
func ConfigUser() {

	interfaces, err := net.Interfaces()
	if err != nil {
		pterm.Error.Println("Не удалось получить настройки сети")
		return
	}

	for _, interf := range interfaces {
		addrs, err := interf.Addrs()
		if err != nil {
			pterm.Error.Println("Ошибка при получении настроек")
			return
		}
		pterm.FgLightBlue.Printf("Сетевой интерфейс: %s\n", interf.Name)

		for _, add := range addrs {
			if ip, ok := add.(*net.IPNet); ok {
				pterm.FgGreen.Printf("\tIP: %v\n", ip)
			}
		}
	}
	port := 2121

	if IsPortOpen(port) {
		pterm.Warning.Println("Вы не можете установить сервер так как порт 2121 уже используется.")

	} else {
		pterm.Success.Println("Порт 2121 не используется! Вы можете запусть сервер")
	}

}
func IsPortOpen(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
