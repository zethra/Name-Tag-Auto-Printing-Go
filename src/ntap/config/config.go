package config
import "log"

type Config struct {
	PrintersFile,
	QueueFile,
	ImagesDirectory,
	ScadDirectory,
	StlDirectory,
	GcodeDirectory,
	DefaultConfig string
	Debug bool
}

func (config *Config) DebugLog(i interface{}) {
	if(config.Debug) {
		log.Println(i)
	}
}