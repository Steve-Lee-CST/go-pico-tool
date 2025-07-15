package request_id

import "github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"

type Config struct {
	HeaderKey         string
	IDGeneratorConfig id_generator.Config
}

var defaultConfig = Config{
	HeaderKey:         "X-Request-ID",
	IDGeneratorConfig: id_generator.GetDefaultConfig(),
}

func GetDefaultConfig() Config {
	return defaultConfig
}
