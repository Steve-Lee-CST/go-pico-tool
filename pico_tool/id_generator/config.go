package id_generator

import (
	"fmt"
	"time"

	"github.com/Steve-Lee-CST/go-pico-tool/pico_tool/utils"
)

type Modifier func(
	timestamp int64, microSecond int64, randSegment string,
) []string

type Config struct {
	Separator *string
	Modifier  Modifier
}

var defaultConfig = Config{
	Separator: utils.ToPtr("-"),
	Modifier: func(timestamp int64, microSecond int64, randSegment string) []string {
		return []string{
			time.Unix(timestamp, 0).Format("20060102150405"),
			fmt.Sprintf("%d", microSecond),
			randSegment,
		}
	},
}

func GetDefaultConfig() Config {
	return defaultConfig
}
