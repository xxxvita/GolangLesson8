package configurator

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// export MYAPP_PORT=8080
// export DB_URL=github.com/xxxvita
type Struct struct {
	PortDefault int    `envconfig:"MYAPP_PORT" default:"8081"`
	DB_URL      string `envconfig:"DB_URL" default:"gb.ru"`
}

// Если указаны файлы, то ими дополняться ENV-переменные, если нет, то загрузятся
// и перетрут существующие из .conf-файлов (эксперементирую)
func Load(envFile ...string) (*Struct, error) {
	var s = Struct{}

	if len(envFile) != 0 {
		godotenv.Load(envFile...)
	} else {
		// Загрузка ENV по умолчанию
		godotenv.Overload("global.conf", "global.conf.private")
	}

	err := envconfig.Process("myapp", &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
