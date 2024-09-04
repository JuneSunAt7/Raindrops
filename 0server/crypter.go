package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const keyFileName = "encryption_key.txt"

// Функция для получения ключа из файла
func getKey() ([]byte, error) {
	key, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(key), nil
}

// Функция для сохранения ключа в файл
func SaveKey(key []byte) error {
	return ioutil.WriteFile(keyFileName, key, 0600)
}

// Функция шифрования
func encrypt(data []byte) ([]byte, error) {
	key, err := getKey()
	if err != nil {
		return nil, fmt.Errorf("cannot get key: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Создаём IV (инициализационный вектор)
	iv := bytes.Repeat([]byte{0}, aes.BlockSize) // Используем нулевой вектор для простоты
	mode := cipher.NewCBCEncrypter(block, iv)

	// Выравниваем длину данных
	data = pad(data)
	ciphertext := make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)
	return ciphertext, nil
}

// Функция выравнивания данных
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Функция шифрования файлов в папке
func EncryptFilesInDir(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(directory, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			encryptedContent, err := encrypt(content)
			if err != nil {
				return err
			}

			// Записываем зашифрованный файл
			encFilePath := filePath + ".enc"
			err = ioutil.WriteFile(encFilePath, encryptedContent, 0644)
			if err != nil {
				return err
			}

			// Удаляем незашифрованный файл
			err = os.Remove(filePath)
			if err != nil {
				return err
			}

			fmt.Printf("Зашифрован и удалён файл: %s", filePath)
		}
	}
	return nil
}
