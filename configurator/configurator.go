package configurator

import (
	"flag"
	"log"
	"os"
	"strings"

	"encoding/json"
	_ "encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// export MYAPP_PORT=8080
// export DB_URL=github.com/xxxvita
type WebConfig struct {
	PortDefault int    `envconfig:"MYAPP_PORT" default:"8081" json:"PortDefault"`
	DbURL       string `envconfig:"DB_URL" default:"gb.ru" json:"DBURL"`
}

var (
	flagPortDefault = flag.Int("myapp-port", 80, "Порт сервиса по умолчанию")
	flagDbURL       = flag.String("db-url", "xyzet.ru", "Порт сервиса по умолчанию")
)

// Если указаны файлы, то ими дополняться ENV-переменные, если нет, то загрузятся
// и перетрут существующие из .conf-файлов (эксперементирую)
func Load(envFile ...string) (*WebConfig, error) {
	var s = WebConfig{}
	filesConfigJSON := make([]string, 0)
	filesConfig := make([]string, 0)

	// Наверное самый простой способ использования флагов впместе с переменными окружения
	// это использования чего-то одного. Если есть флаги - не используем переменные окружения.
	// Для более серьёзной обработки можно, наверное переключиться на cobra или проверять
	// заполненные флаги в командной строке самому.
	if len(os.Args) > 1 {
		flag.Parse()

		if flagPortDefault != nil {
			s.PortDefault = *flagPortDefault
		}

		if flagDbURL != nil {
			s.DbURL = *flagDbURL
		}
	} else {
		// Загрузка из json-файла
		for _, fileName := range envFile {
			if strings.Contains(fileName, ".json") {
				filesConfigJSON = append(filesConfigJSON, fileName)
			} else {
				filesConfig = append(filesConfig, fileName)
			}
		}
	}

	if len(filesConfig) != 0 {
		err := LoadConfig(&s, filesConfig)
		if err != nil {
			return nil, err
		}
	}

	if len(filesConfigJSON) != 0 {
		err := LoadConfigJSON(&s, filesConfigJSON)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func LoadConfig(webConfig *WebConfig, files []string) error {
	if len(files) != 0 {
		godotenv.Load(files...)
	} else {
		// Загрузка ENV по умолчанию
		godotenv.Overload("global.conf", "global.conf.private")
	}

	err := envconfig.Process("myapp", webConfig)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfigJSON(webConfig *WebConfig, files []string) error {
	for _, f := range files {
		fConfig, err := os.Open(f)
		if err != nil {
			log.Fatalf("Ошибка открытия файла-конфигурации в формате JSON: %s", f)
		}

		defer func() {
			err := fConfig.Close()
			if err != nil {
				log.Fatalf("Ошибка закрытия файла-конфигурации в формате JSON: %s", f)
			}

			return
		}()

		err = json.NewDecoder(fConfig).Decode(webConfig)
		if err != nil {
			log.Fatalf("Ошибка парсинга JSON-файла: %s (%s)", f, err)
		}

	}

	return nil
}
