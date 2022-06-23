package setting

import (
	"encoding/json"
	"log"
	"os"
)

type SettingJSON struct {
	Mysql struct {
		Database string `json:"database"`
		User     string `json: "user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json: "port"`
	}
	Mongodb struct {
		Database string `json:"database"`
		User     string `json: "user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json: "port"`
	}
	Nats struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	Btb_service struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	User_service struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
}

func LoadConfigSettingJSON() (SettingJSON, error) {
	var config SettingJSON

	file, err := os.Open("setting/setting.json")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return config, err
}
