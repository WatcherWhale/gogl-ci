package token

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func SaveToken(instance string, token string) error {
	store, err := getStore()
	if err != nil {
		store = map[string]string{}
	}

	store[instance] = token

	jsonStore, err := json.Marshal(store)
	if err != nil {
		return err
	}

	p, err := getPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(p, jsonStore, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GetToken(instance string) (string, error) {
	store, err := getStore()
	if err != nil {
		return "", err
	}

	token, ok := store[instance]
	if !ok {
		return "", fmt.Errorf("%s is not logged in, login with `gogl login %s <token>`", instance, instance)
	}

	return token, nil
}

func getStore() (map[string]string, error) {
	var store map[string]string

	p, err := getPath()

	if err != nil {
		return store, err
	}

	storeJson, err := os.ReadFile(p)
	if err != nil {
		return store, err
	}

	err = json.Unmarshal(storeJson, &store)
	if err != nil {
		return store, err
	}

	return store, nil
}

func getPath() (string, error) {
	configHome, err := os.UserConfigDir()

	if err != nil {
		return "", nil
	}

	configPath := path.Join(configHome, "gogl")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	tokenStore := path.Join(configPath, "login.json")

	return tokenStore, nil
}
