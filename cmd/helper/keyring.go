package helper

import (
	"os"

	"github.com/zalando/go-keyring"
)

const keyringSvc = "dbac"

func noKeychain() bool {
	return os.Getenv("DBAC_NO_KEYCHAIN") != ""
}

func StorePassword(profileName, password string) error {
	if noKeychain() {
		return storeAESPassword(profileName, password)
	}
	return keyring.Set(keyringSvc, profileName, password)
}

func RetrievePassword(profileName string) (string, error) {
	if noKeychain() {
		return retrieveAESPassword(profileName)
	}
	return keyring.Get(keyringSvc, profileName)
}

func DeletePassword(profileName string) error {
	if noKeychain() {
		return deleteAESPassword(profileName)
	}
	return keyring.Delete(keyringSvc, profileName)
}
