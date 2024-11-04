package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"errors"

)

func CBCDecrypter( ciphertext []byte) ([]byte, error) {
	key := md5.Sum([]byte(SESSION_PASSWD))

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