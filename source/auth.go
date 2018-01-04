package source

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	log "github.com/Sirupsen/logrus"
)

type novaConfigScheme struct {
	NovaURL    string `json:"NOVA_URL"`
	NovaID     string `json:"NOVA_CLIENT_ID"`
	NovaSecret string `json:"NOVA_CLIENT_SECRET"`
}

type configFile struct {
	Nova []novaConfigScheme
}

func GetBasicAuthHeader(clientID, clientSecret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
}

func GetCredentials(novaUrl string) (clientID, clientSecret string, err error) {
	errNotFound := fmt.Errorf("couldn't find credentials in ENV[NOVA_CLIENT_ID], ENV[NOVA_CLIENT_SECRET] or in config file")

	// env takes precedence over file vars
	clientID = os.Getenv("NOVA_CLIENT_ID")
	clientSecret = os.Getenv("NOVA_CLIENT_SECRET")
	if clientID != "" && clientSecret != "" {
		return
	}

	// try file vars if env is empty
	rawFile, err := ioutil.ReadFile(getConfigFilePath())
	if err != nil {
		log.Errorf("Error reading %s: %+v", getConfigFilePath(), err)
		err = errNotFound
		return
	}
	novaConfig := configFile{}
	json.Unmarshal(rawFile, &novaConfig)
	for _, nc := range novaConfig.Nova {
		if nc.NovaURL == novaUrl {
			clientID, clientSecret = nc.NovaID, nc.NovaSecret
		}
	}
	if clientID != "" && clientSecret != "" {
		return
	}
	log.Errorf("couldn't find %s in %s", novaUrl, getConfigFilePath())
	err = errNotFound
	return
}

func SaveCredentials(novaUrl string) (clientID, clientSecret string, err error) {
	var n int

	currentConfig := configFile{}
	rawFile, _ := ioutil.ReadFile(getConfigFilePath())
	if rawFile != nil {
		json.Unmarshal(rawFile, &currentConfig)
		for _, nc := range currentConfig.Nova {
			if nc.NovaURL == novaUrl {
				err = fmt.Errorf("%s already contains an entry for %s. Please modify the file by hand to update keys.",
					getConfigFilePath(), novaUrl)
				return
			}
		}
	}

	fmt.Println("You can get your splunknova client credentials at https://www.splunknova.com/apikeys\n")

	fmt.Printf("Please enter Client ID: ")
	n, err = fmt.Scan(&clientID)
	if n != 1 || err != nil {
		err = fmt.Errorf("error reading Client ID %+v", err)
		return
	}

	fmt.Printf("Please enter Client Secret: ")
	n, err = fmt.Scan(&clientSecret)
	if n != 1 || err != nil {
		err = fmt.Errorf("error reading Client Secret %+v", err)
		return
	}

	err = validateCredentials(novaUrl, clientID, clientSecret)
	if err != nil {
		return
	}
	log.Infof("Login succeeded")

	currentConfig.Nova = append(currentConfig.Nova, novaConfigScheme{novaUrl, clientID, clientSecret})

	data, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(getConfigFilePath(), data, 0644)
	if err != nil {
		log.Infof("Keys saved to %s", getConfigFilePath())
	}
	return
}

func validateCredentials(novaUrl, clientID, clientSecret string) error {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	_, err := Get(novaUrl+validateCredsURLPath, nil, authHeader)
	if err != nil {
		return fmt.Errorf("credentials didn't work, please try again or contact us")
	}
	return nil
}

func getConfigFilePath() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return dir + configFileRelPath
}
