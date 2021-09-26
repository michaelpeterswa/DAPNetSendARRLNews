package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var settings DapnetNewsSettings

func main() {
	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &settings)
	if err != nil {
		logger.Error("YAML failed to unmarshal to DapnetSettings", zap.Error(err))
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(settings.NewsEndpoint)
	if err != nil {
		logger.Error("Couldn't create feed object from URL", zap.Error(err))
	}
	for _, article := range feed.Items {
		str := fmt.Sprintf("%s:%s", article.Title, article.PublishedParsed)
		str = strings.ReplaceAll(str, " ", "")
		hash := sha256.Sum256([]byte(str))
		hex := hex.EncodeToString(hash[:])
		fmt.Println(hex)
	}
}
