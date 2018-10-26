package configuration

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
	"os"
	"path/filepath"
)

type Configuration interface {
	GetDatabase() Database
	GetKeys() Keys
	GetMail() Mails
	GetUrl() string
	loadConfiguration()
	GetFilePath(filename string) string
}

type FileConfiguration struct {
	Db   Database
	Key  Keys
	Mail Mails
	Url  string
}

func (conf *FileConfiguration) loadConfiguration() {
	buffer, err := ioutil.ReadFile(conf.GetFilePath(CONF_FILENAME))
	if err != nil {
		log.Fatalf("Erreur dans le chargement de la configuration.", err.Error())
	}
	goyaml.Unmarshal(buffer, conf)
}

func (conf *FileConfiguration) GetFilePath(filename string) string {
	folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Erreur dans le path du fichier. %s", err.Error())
	}
	return folder + "/" + filename
}

func (conf *FileConfiguration) GetDatabase() Database {
	return conf.Db
}

func (conf *FileConfiguration) GetKeys() Keys {
	return conf.Key
}

func (conf *FileConfiguration) GetMail() Mails {
	return conf.Mail
}

func (conf *FileConfiguration) GetUrl() string {
	return conf.Url
}

type TestConfiguration struct{}

func (conf *TestConfiguration) loadConfiguration() {}

func (conf *TestConfiguration) GetFilePath(filename string) string {
	return filename
}

func (conf *TestConfiguration) GetDatabase() Database {
	return Database{"localhost", "reunion", "test", "test"}
}

func (conf *TestConfiguration) GetMail() Mails {
	return Mails{"test@test", "contact@animalove.re"}
}

func (conf *TestConfiguration) GetUrl() string {
	return "localhost:8080"
}

func (conf *TestConfiguration) GetKeys() Keys {
	return Keys{Private: "../token/private_key", Public: "../token/public_key"}
}
