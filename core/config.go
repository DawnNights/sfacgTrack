package core

import (
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