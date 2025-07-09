package id_generator

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type IDGenerator struct {
	config Config
}

func NewIDGenerator(config Config) *IDGenerator {
	return &IDGenerator{
		config: config,
	}
}

func (tool *IDGenerator) Generate() string {
	timestamp, microSeconds := func() (int64, int64) {
		now := time.Now()
		return now.Unix(), now.UnixMicro() % 1e6
	}()
	randSegment := strings.Split(uuid.NewString(), "-")[0]

	modifier, separator := tool.config.Modifier, tool.config.Separator
	if modifier == nil {
		modifier = GetDefaultConfig().Modifier
	}
	if separator == nil {
		separator = GetDefaultConfig().Separator
	}
	return strings.Join(modifier(timestamp, microSeconds, randSegment), *separator)
}
