package server

import (
	"strings"
	"encoding/binary"
	"io"
	"log"
	"fmt"
	"net"
	"time"
	"os"
	"io/ioutil"
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
    repoPath := filepath.Join(ROOT, Uname, "repositories", nameOfrepo) 
    log.Printf("Создание каталога: %s\n", repoPath)
    log.Printf("ROOT: %s, Uname: %s, nameOfrepo: %s", ROOT, Uname, nameOfrepo)

    defer conn.Close()

    if err := os.MkdirAll(repoPath, 0777); err != nil {
        log.Println("Ошибка при создании каталога:", err)
        return
    }


    for {
		conn.Write([]byte("200 Start upload!\n"))
       
        fileNameBuf := make([]byte, 256) 
        n, err := conn.Read(fileNameBuf)
        if err != nil {
            log.Println("Ошибка при чтении имени файла:", err)
            break
        }

        fileName := strings.Trim(string(fileNameBuf[:n]), "\x00\n") 
        if fileName == "" {
            log.Println("Получено пустое имя файла, завершение...")
            break 
        }

       
        var fileSize int64
        err = binary.Read(conn, binary.BigEndian, &fileSize)
        if err != nil {
            log.Println("Ошибка при чтении размера файла:", err)
            break
        }

        
        filePath := filepath.Join(repoPath, fileName) 
        outFile, err := os.Create(filePath)
        if err != nil {
            log.Println("Ошибка при создании файла:", err)
            break
        }

       
        _, err = io.CopyN(outFile, conn, fileSize)
        if err != nil {
            log.Println("Ошибка при копировании содержимого файла:", err)
            outFile.Close() 
            break
        }

        outFile.Close() 
        log.Printf("Файл %s успешно сохранён (%d байт).\n", fileName, fileSize)
    }

    log.Println("Передача файлов завершена.")
}
func GetListKeeps(conn net.Conn, nameOfrepo string) {
    files, err := ioutil.ReadDir(ROOT + "/" + Uname + "/repositories/" + nameOfrepo) // Read all filenames from filestore.
	if err != nil {
		pterm.Error.Println("No keeps in current repository")	
		conn.Write([]byte("No keeps in current repository"))
		log.Println(err.Error())
	}

	fileINFO := ""
	for _, file := range files {
		if !file.IsDir() {
			fileINFO += fmt.Sprintf("%-40s%-25s%-10d\n",
				file.Name(),
				file.ModTime().Format("2006-01-02 15:04:05"),
				file.Size())
		}

	}
	conn.Write([]byte(fileINFO))
}
func GetLastKeep(conn net.Conn, nameOfrepo string) (string, error) {
    dir := filepath.Join(ROOT, Uname, "repositories", nameOfrepo)
	log.Printf("Reading directory: %s\n", dir)
    files, err := os.ReadDir(dir)
    if err != nil {
        return "", fmt.Errorf("error reading directory: %v", err)
    }

    var latestFile os.DirEntry
    var latestTime time.Time
    for _, file := range files {

        if !file.IsDir() {
            info, err := file.Info()
            if err != nil {
                continue
            }

            if latestFile == nil || info.ModTime().After(latestTime) {
                latestFile = file
                latestTime = info.ModTime()
            }
        }
    }
	log.Printf("Latest file: %s\n", latestFile.Name())


    if latestFile == nil {
        return "", fmt.Errorf("no files found in directory: %s", dir)
    }
	log.Printf("Latest file: %s\n", latestFile.Name())

    filePath := filepath.Join(dir, latestFile.Name())
    log.Printf("Sending file: %s\n", filePath)

    fileInfo, err := os.Stat(filePath)
    if err != nil {
        log.Printf("Error getting info about file %s: %v", latestFile.Name(), err)
        return "", fmt.Errorf("error getting info about file %s: %v", latestFile.Name(), err)
    }

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        log.Printf("Error reading from connection: %v", err)
        return "", fmt.Errorf("error reading from connection: %v", err)
    }

    str := strings.TrimSpace(string(buf[:n]))
    commandArr := strings.Fields(str)
    if len(commandArr) == 0 || commandArr[0] != "200" {
        log.Println("Error: expected command 200")
        return "", fmt.Errorf("error signal received from client")
    }

    _, err = conn.Write([]byte(fmt.Sprintf("%s\x00\n", latestFile.Name()))) 
    if err != nil {
        log.Printf("Error sending file name: %v", err)
        return "", err
    }  


    size := fileInfo.Size()
    err = binary.Write(conn, binary.BigEndian, size)
    if err != nil {
        log.Printf("Error sending file size: %v", err)
        return "", err
    }


    f, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening file %s: %v", latestFile.Name(), err)
        return "", fmt.Errorf("error opening file %s: %v", latestFile.Name(), err)
    }
    defer f.Close()

  
    if _, err = io.Copy(conn, f); err != nil {
        log.Printf("Error transferring file %s: %v", latestFile.Name(), err)
        return "", fmt.Errorf("error transferring file %s: %v", latestFile.Name(), err)
    }

    log.Printf("File %s successfully sent.", latestFile.Name())


    _, err = conn.Write([]byte("\x00")) 
    if err != nil {
        log.Printf("Error sending end signal: %v", err)
        return "", err
    }

    return latestFile.Name(), nil
}