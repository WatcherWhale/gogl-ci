package cache

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func CreateCacheKey(filename string) string {
	absPath, err := filepath.Abs(filename)

	if err != nil {
		return base64.StdEncoding.EncodeToString([]byte(filename))
	}

	return base64.StdEncoding.EncodeToString([]byte(absPath))
}

func SaveCache(cacheKey string, obj any) error {
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	cacheDir := path.Join(cacheHome, "gogl")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer
	err = gob.NewEncoder(&b).Encode(obj)

	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(cacheDir, cacheKey), b.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func LoadCache(cacheKey string, obj any) error {
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	cacheDir := path.Join(cacheHome, "gogl")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("cache does not exist")
	}

	b, err := os.ReadFile(path.Join(cacheDir, cacheKey))
	if err != nil {
		return err
	}

	err = gob.NewDecoder(bytes.NewReader(b)).Decode(obj)
	if err != nil {
		return err
	}

	return nil
}

func CleanCache() error {
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	cacheDir := path.Join(cacheHome, "gogl")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil
	}

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(path.Join(cacheDir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
