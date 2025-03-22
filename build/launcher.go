package main

import (
    "fmt"
	"os"
	"path/filepath"
	"io"
    "os/exec"
    "time"
)
func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return err
    }

    return nil
}

func moveFilesFromBin(binDir string, filesToMove []string) error {

    parentDir := filepath.Dir(binDir)

    for _, fileName := range filesToMove {
        srcPath := filepath.Join(binDir, fileName)
        dstPath := filepath.Join(parentDir, fileName)

        if _, err := os.Stat(srcPath); os.IsNotExist(err) {
            continue
        }

        err := copyFile(srcPath, dstPath)
        if err != nil {
            continue
        }

    }

    return nil
}
func main() {
    programs := []string{
        "bin/generate_ssl.exe",                          
    }

    fmt.Println("Setup certificates...")
    for _, program := range programs {
        err := exec.Command(program).Start()
        if err != nil {
            fmt.Printf("Error %s: %v\n", program, err)
    }

    time.Sleep(5 * time.Second)
} 
	fmt.Println("Setup data...")
	fileName := "user_creds.db"

    file, err := os.Create(fileName)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer file.Close() 
	dirName := "filestore"

    err = os.Mkdir(dirName, 0755) 
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
	fmt.Println("Setup apps...")
	binDir := "./bin"

    filesToMove := []string{
        "server.exe",
        "server_mgmt.exe",
    }

    err = moveFilesFromBin(binDir, filesToMove)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

}