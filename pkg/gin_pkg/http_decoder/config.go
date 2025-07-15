package http_decoder

type Config struct {
	HttpRequestKey        string
	HttpResponseKey       string
	HttpResponseWriterKey string
}

var defaultConfig = Config{
	HttpRequestKey:        "http_request",
	HttpResponseKey:       "http_response",
	HttpResponseWriterKey: "http_response_writer",
}

func DefaultConfig() Config {
	return defaultConfig
}
