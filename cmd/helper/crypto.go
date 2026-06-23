package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/denisbrodbeck/machineid"
)

const (
	aesPrefix     = "enc:v1:"
	passwordsFile = "/.dbac-passwords.json"
)

// deriveKey produces a 32-byte AES key from a machine-scoped HMAC keyed to "dbac".
func deriveKey() ([]byte, error) {
	id, err := machineid.ProtectedID("dbac")
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %v", err)
	}
	key := sha256.Sum256([]byte(id))
	return key[:], nil
}

func encryptAES(plain string) (string, error) {
	key, err := deriveKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return aesPrefix + base64.StdEncoding.EncodeToString(ct), nil
}

func decryptAES(s string) (string, error) {
	if !strings.HasPrefix(s, aesPrefix) {
		// plaintext passthrough for legacy migration
		return s, nil
	}
	key, err := deriveKey()
	if err != nil {
		return "", err
	}
	ct, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, aesPrefix))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(ct) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ct := ct[:gcm.NonceSize()], ct[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func aesPasswordsPath() string {
	return os.Getenv("HOME") + passwordsFile
}

func readPasswordsFile() (map[string]string, error) {
	passwords := make(map[string]string)
	data, err := os.ReadFile(aesPasswordsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return passwords, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &passwords); err != nil {
		return nil, err
	}
	return passwords, nil
}

func writePasswordsFile(passwords map[string]string) error {
	data, err := json.MarshalIndent(passwords, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(aesPasswordsPath(), data, 0600)
}

func storeAESPassword(profileName, password string) error {
	encrypted, err := encryptAES(password)
	if err != nil {
		return err
	}
	passwords, err := readPasswordsFile()
	if err != nil {
		return err
	}
	passwords[profileName] = encrypted
	return writePasswordsFile(passwords)
}

func retrieveAESPassword(profileName string) (string, error) {
	passwords, err := readPasswordsFile()
	if err != nil {
		return "", err
	}
	encrypted, ok := passwords[profileName]
	if !ok {
		return "", fmt.Errorf("password not found for profile '%s'", profileName)
	}
	return decryptAES(encrypted)
}

func deleteAESPassword(profileName string) error {
	passwords, err := readPasswordsFile()
	if err != nil {
		return err
	}
	delete(passwords, profileName)
	if len(passwords) == 0 {
		return os.Remove(aesPasswordsPath())
	}
	return writePasswordsFile(passwords)
}
