package utils

import (
	"encoding/json"
	"github.com/Ypxd/diplom/auth/internal/models"
	"log"
	"os"
	"sync"
)

const confFileName = "configurations/config.json"

var (
	once   sync.Once
	config *models.Config
)

func GetConfig() *models.Config {
	once.Do(func() {
		config = &models.Config{}
		f, err := os.Open(confFileName)
		if err != nil {
			log.Fatalln(err)
			return
		}
		fi, err := f.Stat()
		if err != nil {
			log.Fatalln(err)
			return
		}
		data := make([]byte, fi.Size())
		_, err = f.Read(data)
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = json.Unmarshal(data, config)
		if err != nil {
			log.Fatalln(err)
		}
	})
	if config == nil {
		config = &models.Config{
			Server: models.Server{
				Host: "localhost",
				Port: 9095,
			},
		}
	}

	return config
}
