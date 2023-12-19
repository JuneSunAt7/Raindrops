package client

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

func Calendar() {
	file, err := os.Open(ROOT + "/" + "localSettings" + "/" + "settings.ini")
	if err != nil {
		pterm.FgLightRed.Println("Файл настроек не найден!")
		return
	}
	defer file.Close()

	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		pterm.BgCyan.Println("Дни резервирования:")

		pterm.FgGreen.Println(strings.ReplaceAll(string(data[:n]), " ", "\n"))

	}
}

func Userfiles() {

	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Println("Доступные директории:")

	var options []string
	maindir := ChooseDir()
	files, err := ioutil.ReadDir(maindir)
	if err != nil {
		pterm.FgRed.Println("Ошибка чтения директорий и файлов!")
	}

	for _, file := range files {
		absPath, err := filepath.Abs(maindir + "\\" + file.Name())
		if err != nil {
			pterm.FgRed.Println("Ошибка прочтения пути к файлу!")
		}
		options = append(options, fmt.Sprint(absPath+"\n"))

	}
	updateSettings(options)
	pterm.Success.Println("Сохранено в настройках")
}

func updateSettings(files []string) {
	// Write path selected files
	outputFile, err := os.Create(ROOT + "/" + "localSettings" + "/" + "path.ini")
	if err != nil {
		pterm.FgRed.Printfln("Ошибка создания локального файла!")
	}
	defer outputFile.Close()
	outputFile.WriteString(strings.Join(files, "\n"))
}
func Setting() {

	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Println("Дни для резервирования")

	var options []string
	options = append(options, fmt.Sprintf("Понедельник"))
	options = append(options, fmt.Sprintf("Вторник"))
	options = append(options, fmt.Sprintf("Среда"))
	options = append(options, fmt.Sprintf("Четверг"))
	options = append(options, fmt.Sprintf("Пятница"))
	options = append(options, fmt.Sprintf("Суббота"))
	options = append(options, fmt.Sprintf("Воскресенье"))

	selectedOptions, _ := pterm.DefaultInteractiveMultiselect.WithOptions(options).Show()
	pterm.Info.Printfln("Выбранные дни для резервирования: %s", pterm.Green(selectedOptions))
	createSettingsFile(selectedOptions)
}

func createSettingsFile(days []string) {
	/* errmk := os.Mkdir(ROOT+"/"+"localSettings", 0777)
	if errmk != nil {

		fmt.Println(errmk)
		pterm.FgLightRed.Println("Ошибка создания реестра настроек!")
	} */

	outputFile, err := os.Create(ROOT + "/" + "localSettings" + "/" + "settings.ini")
	if err != nil {
		pterm.FgRed.Printfln("Ошибка создания локального файла!")
	}
	defer outputFile.Close()
	outputFile.WriteString(strings.Join(days, "\n"))

}

func Containers() {
	folderPath := ChooseDir()
	folderName := filepath.Base(folderPath)

	// Создаем новый zip архив для записи
	zipFile, err := os.Create(folderName + ".rdct")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer zipFile.Close()

	// Создаем новый zip писатель
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Рекурсивно обходим все файлы и подпапки в указанной папке
	err = filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Создание контейнера...").Start()

		for i := 0; i < p.Total; i++ {
			// Progressbae - uploader
			p.UpdateTitle("Создание контейнера")
			p.Increment()
		}

		if err != nil {
			return err
		}

		// Игнорируем директории
		if fileInfo.IsDir() {
			return nil
		}

		// Относительный путь файла внутри папки
		relativePath := strings.TrimPrefix(filePath, folderPath)

		// Создаем заголовок файла в архиве
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		// Устанавливаем имя файла в архиве
		header.Name = relativePath

		// Устанавливаем метод сжатия
		header.Method = zip.Deflate

		// Создаем запись в архиве
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Открываем существующий файл для чтения
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Копируем содержимое файла в запись архива
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	pterm.Success.Println("Успешная контейнеризация!")

}
func dirs(folderPath string) {
	
	folderName := filepath.Base(folderPath)

	// Создаем новый zip архив для записи
	zipFile, err := os.Create(folderName + ".rdct")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer zipFile.Close()

	// Создаем новый zip писатель
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Рекурсивно обходим все файлы и подпапки в указанной папке
	err = filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Создание контейнера...").Start()

		for i := 0; i < p.Total; i++ {
			// Progressbae - uploader
			p.UpdateTitle("Создание контейнера")
			p.Increment()
		}

		if err != nil {
			return err
		}

		// Игнорируем директории
		if fileInfo.IsDir() {
			return nil
		}

		// Относительный путь файла внутри папки
		relativePath := strings.TrimPrefix(filePath, folderPath)

		// Создаем заголовок файла в архиве
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		// Устанавливаем имя файла в архиве
		header.Name = relativePath

		// Устанавливаем метод сжатия
		header.Method = zip.Deflate

		// Создаем запись в архиве
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Открываем существующий файл для чтения
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Копируем содержимое файла в запись архива
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	pterm.Success.Println("Успешная контейнеризация!")

}
