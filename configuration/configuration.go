//author Julian Klempa (4085242)

package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Credentials map[string][]byte
	SavesCount  int
	Port        string
}

func ReadConfig() Configuration {
	file, _ := os.OpenFile("./configuration/config.json", os.O_RDONLY|os.O_CREATE, 0777)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}

func WriteConfig(configuration Configuration) {
	file, _ := os.OpenFile("./configuration/config.json", os.O_WRONLY|os.O_CREATE, 0777)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ")
	encoder.Encode(configuration)
}

func GetCredentials() map[string][]byte {
	conf := ReadConfig()
	return conf.Credentials
}

func SetConfig(credentials map[string][]byte, savesCount int, port string) {
	if savesCount > 20 {
		savesCount = 20
	} else if savesCount < 0 {
		savesCount = 0
	}

	conf := Configuration{Credentials: credentials, SavesCount: savesCount, Port: port}
	WriteConfig(conf)
}
