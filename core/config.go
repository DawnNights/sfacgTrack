package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

func LoadConfig() TrackConfig {
	var sf SFAPI
	var config TrackConfig

	yaml.Unmarshal(FileRead("resource\\Config.yaml"),&config)
	for idx, _ := range config {
		config[idx].RecordUrl = sf.FindChapterUrl(config[idx].BookId)
	}

	return config
}

func (config TrackConfig) Print() {
	for idx, elem := range config {
		fmt.Println("*",idx+1,elem.BookId,elem.RecordUrl)
	}
}