package server

import (
	"io"
	"net"
	"os"
	"log"
	"bytes"
	"path/filepath"
	"encoding/binary"

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
    // Открываем директорию для сохранения
    repoPath := ROOT + "/" + Uname + "/repositories/" + nameOfrepo
    if _, err := os.Stat(repoPath); os.IsNotExist(err) {
        os.MkdirAll(repoPath, 0777)
    }

    conn.Write([]byte("200 Start upload!"))

    for {
        // Читаем имя файла
        fileNameBuf := make([]byte, 50) // Максимальная длина имени файла
        _, err := io.ReadFull(conn, fileNameBuf)
        if err != nil {
            log.Println("Ошибка при получении имени файла:", err)
            break
        }

        fileName := string(bytes.Trim(fileNameBuf, "\x00")) // Удаляем возможные нулевые байты
        if fileName == "" {
            break // Если имя пустое, выходим из цикла
        }

        // Читаем размер файла
        var fileSize int64
        err = binary.Read(conn, binary.BigEndian, &fileSize)
        if err != nil {
            log.Println("Ошибка при получении размера файла:", err)
            break
        }

        // Создаем файл
        filePath := filepath.Join(repoPath, fileName)
        out, err := os.Create(filePath)
        if err != nil {
            log.Println("Ошибка при создании файла:", err)
            break
        }
        defer out.Close() // Закрываем файл после завершения работы

        // Копируем данные файла
        _, err = io.CopyN(out, conn, fileSize)
        if err != nil {
            log.Println("Ошибка при сохранении файла:", err)
            break
        }
        log.Printf("Файл %s успешно загружен.\n", fileName)
    }
    conn.Write([]byte("Upload complete!"))
}