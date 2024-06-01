package config

var EnvCfg envConfig

type envConfig struct {
	LoggerLevel          string  `env:"LOGGER_LEVEL" envDefault:"INFO"`
	LoggerWithTraceState string  `env:"LOGGER_OUT_TRACING" envDefault:"disable"`
	TiedLogging          string  `env:"TIED" envDefault:"NONE"`
	TracingEndPoint      string  `env:"TRACING_ENDPOINT"`
	PyroscopeAddr        string  `env:"PYROSCOPE_ADDR"`
	PyroscopeState       string  `env:"PYROSCOPE_STATE" envDefault:"false"`
	PodIpAddr            string  `env:"POD_IP" envDefault:"127.0.0.1"`
	OtelState            string  `env:"TRACING_STATE" envDefault:"enable"`
	OtelSampler          float64 `env:"TRACING_SAMPLER" envDefault:"0.01"`
}
