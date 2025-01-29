package codegen

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"strings"

	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
)

const (
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 8
)

func Generate(repos *repository.Repository) error {
	used, err := repos.Generator.CountUsedCodes()
	if err != nil {
		return err
	}

	if used > 3000 || used == 0 {
		secreKey := []byte(os.Getenv("SYMMETRICK_KEY"))

		adjectives, err := loadWords("./pkg/codegen/adjectives.txt")
		if err != nil {
			return err
		}
		nouns, err := loadWords("./pkg/codegen/nouns.txt")
		if err != nil {
			return err
		}
		for i := 0; i < 3000; i++ {
			codeWord := generateCodeWord(adjectives, nouns)
			password := generateSimplePassword()
			password, err = Encrypt(password, secreKey)
			if err != nil {
				return err
			}

			err = repos.Generator.SaveToDB(codeWord, password)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateCodeWord(adjectives, nouns []string) string {
	adj := adjectives[mrand.Intn(len(adjectives))]
	noun := nouns[mrand.Intn(len(nouns))]
	number := mrand.Intn(100)

	return fmt.Sprintf("%s_%s%d", adj, noun, number)
}

func generateSimplePassword() string {
	var passwordBuilder strings.Builder
	passwordBuilder.Grow(length)

	for i := 0; i < length; i++ {
		num := mrand.Intn(60)
		passwordBuilder.WriteByte(symbols[num])
	}

	return passwordBuilder.String()
}

func Encrypt(text string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string, secretKey []byte) (string, error) {
    ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

    block, err := aes.NewCipher(secretKey)
    if err != nil {
        return "", err
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}

func loadWords(fileName string) ([]string, error) {
	var words []string

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
