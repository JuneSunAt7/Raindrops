package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"math"

	"os"
)

func CBCEncrypter(password string, sl []byte) ([]byte, error) {
	key := md5.Sum([]byte(PASSWD))
	sl16 := make([]byte, int(math.Ceil(float64(len(sl))/aes.BlockSize)*aes.BlockSize)) //%16 bytes
	copy(sl16, sl)

	if len(sl16)%aes.BlockSize != 0 {

		return nil, errors.New("plaintext is not a multiple of the block size")

	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(sl16))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, (err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], sl16)

	return ciphertext, nil
}

func CBCDecrypter(password string, ciphertext []byte) ([]byte, error) {
	key := md5.Sum([]byte(password))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, (err)
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	ciphertextOut := make([]byte, len(ciphertext))
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertextOut, ciphertext)

	ciphertextOut = bytes.TrimRight(ciphertextOut, "\x00")
	ciphertextOut = ciphertextOut[:len(ciphertextOut)-1] //именно 1!
	return ciphertextOut, nil
}

func CheckFileMD5Hash(path string) {

	hashFile, _ := os.Open(path)
	defer hashFile.Close()
	h := md5.New()
	if _, err := io.Copy(h, hashFile); err != nil {
		log.Println("error calc hash:", err)
		return
	}
	log.Printf("Хеш: %x\n", h.Sum(nil))
}
