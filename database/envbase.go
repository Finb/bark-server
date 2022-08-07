package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
)

type EnvBase struct {
}

const (
	MARKDOWN_DIR = "markdown"
)

func NewEnvBase() Database {
	return &EnvBase{}
}

func (d *EnvBase) CountAll() (int, error) {
	return 1, nil
}

func (d *EnvBase) DeviceTokenByKey(key string) (string, error) {
	if key == os.Getenv("BARK_KEY") {
		return os.Getenv("BARK_DEVICE_TOKEN"), nil
	}
	return "nil", fmt.Errorf("key not found")
}

func (d *EnvBase) SaveDeviceTokenByKey(key, token string) (string, error) {
	if token == os.Getenv("BARK_DEVICE_TOKEN") {
		return os.Getenv("BARK_KEY"), nil
	}
	return "nil", fmt.Errorf("device token is invalid")
}

// Get Markdown Content by key
func (d *EnvBase) GetMarkdownByKey(key string) (string, error) {
	mdPath := filepath.Join(MARKDOWN_DIR, key+".md")
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return "", err
	}
	return string(content), err
}

// Save Markdown Content
func (d *EnvBase) SaveMarkdown(content string) (string, error) {
	if _, err := os.Stat(MARKDOWN_DIR); os.IsNotExist(err) {
		logger.Infof("init markdown storage dir [%s]...", MARKDOWN_DIR)
		if err = os.MkdirAll(MARKDOWN_DIR, 0755); err != nil {
			logger.Fatalf("failed to create markdown storage dir(%s): %v", MARKDOWN_DIR, err)
		}
	} else if err != nil {
		logger.Fatalf("failed to open markdown storage dir(%s): %v", MARKDOWN_DIR, err)
	}
	key := shortuuid.New()
	mdPath := filepath.Join(MARKDOWN_DIR, key+".md")
	err := ioutil.WriteFile(mdPath, []byte(content), 0666)
	if err != nil {
		logger.Fatalf("failed to create markdown: %v", err)
	}
	return key, err
}

func (d *EnvBase) Close() error {
	return nil
}
