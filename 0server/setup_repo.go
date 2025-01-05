package server

import (
	"strings"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
)

func SetupRepo(conn net.Conn){
	errMkDir := os.MkdirAll(ROOT+"/"+Uname+"/repositories", 0777)
	if errMkDir != nil {
		pterm.Warning.Println("Ошибка создания папки", errMkDir)
		return
	}
	conn.Write([]byte("Success setup repo"))
}
func GetKeeps(conn net.Conn, nameOfrepo string) {
    repoPath := filepath.Join(ROOT, Uname, "repositories", nameOfrepo) // Используем filepath.Join
    log.Printf("Создание каталога: %s\n", repoPath)
    log.Printf("ROOT: %s, Uname: %s, nameOfrepo: %s", ROOT, Uname, nameOfrepo)

    defer conn.Close()

    // Ответ клиенту о начале загрузки
    //conn.Write([]byte("200 Start upload!\n"))

    // Создаём каталог для сохранения загруженных файлов
    if err := os.MkdirAll(repoPath, 0777); err != nil {
        log.Println("Ошибка при создании каталога:", err)
        return
    }

    for {
		conn.Write([]byte("200 Start upload!\n"))
        // Чтение имени файла
        fileNameBuf := make([]byte, 256) // Размер буфера для имени файла
        n, err := conn.Read(fileNameBuf)
        if err != nil {
            log.Println("Ошибка при чтении имени файла:", err)
            break
        }

        fileName := strings.Trim(string(fileNameBuf[:n]), "\x00\n") // Удаляем нулевые байты и новую строку
        if fileName == "" {
            log.Println("Получено пустое имя файла, завершение...")
            break // Завершаем при получении пустого имени
        }

        // Чтение размера файла
        var fileSize int64
        err = binary.Read(conn, binary.BigEndian, &fileSize)
        if err != nil {
            log.Println("Ошибка при чтении размера файла:", err)
            break
        }

        // Создаём файл для сохранения
        filePath := filepath.Join(repoPath, fileName) // Используем filepath.Join
        outFile, err := os.Create(filePath)
        if err != nil {
            log.Println("Ошибка при создании файла:", err)
            break
        }

        // Копируем содержимое файла
        _, err = io.CopyN(outFile, conn, fileSize)
        if err != nil {
            log.Println("Ошибка при копировании содержимого файла:", err)
            outFile.Close() // Закрываем файл при ошибке
            break
        }

        outFile.Close() // Закрываем файл после завершения копирования
        log.Printf("Файл %s успешно сохранён (%d байт).\n", fileName, fileSize)
    }

    log.Println("Передача файлов завершена.")
}