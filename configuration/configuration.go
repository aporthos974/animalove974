package configuration

const CONF_FILENAME = "configuration.yml"

var Conf Configuration

type Database struct {
	Url, Name, Username, Password string
}

type Keys struct {
	Private, Public string
}

type Mails struct {
	Admin, Contact string
}

func NewInstance() {
	Conf = new(FileConfiguration)
	Conf.loadConfiguration()
}

func GetConfiguration() Configuration {
	return Conf
}
