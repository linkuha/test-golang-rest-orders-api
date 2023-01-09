package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	Version           = "1.0.0beta"
	DefaultLogDir     = "/var/log"
	DefaultConfigPath = "./config/config.yml"
	DefaultLogLevel   = "info"
)

type AppFlags struct {
	LogDir     stringFlag
	LogLevel   stringFlag
	ConfigPath stringFlag
}

type EnvParams struct {
	FileLoaded  bool
	Env         string `mapstructure:"APP_ENV" env:"APP_ENV"`
	LogDir      string `mapstructure:"APP_LOG_DIR" env:"APP_LOG_DIR"`
	LogLevel    string `mapstructure:"APP_LOG_LEVEL" env:"APP_LOG_LEVEL"`
	Port        string `mapstructure:"PORT" env:"PORT"`
	Host        string `mapstructure:"API_HOST" env:"API_HOST"`
	PgHost      string `mapstructure:"POSTGRES_HOST" env:"POSTGRES_HOST"`
	PgPort      string `mapstructure:"POSTGRES_PORT" env:"POSTGRES_PORT"`
	PgUser      string `mapstructure:"POSTGRES_USER" env:"POSTGRES_USER"`
	PgPassword  string `mapstructure:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	PgDB        string `mapstructure:"POSTGRES_DB" env:"POSTGRES_DB"`
	DatabaseURL string `mapstructure:"DATABASE_URL" env:"DATABASE_URL"`
}

type FileParams struct {
	FileLoaded bool
	Test       string `mapstructure:"test_string"`
}

type MergedParams struct {
	LogInFile  bool
	LogDir     string
	LogLevel   string
	ConfigPath string
}

type Config struct {
	EnvParams        EnvParams
	FileParams       FileParams
	ConfigFileLoaded bool
	flags            AppFlags
	Merged           MergedParams
}

type stringFlag struct {
	set   bool
	value string
}

type boolFlag struct {
	set   bool
	value bool
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func (sf *boolFlag) Set(x bool) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *boolFlag) String() string {
	if sf.value == true {
		return "true"
	}
	return "false"
}

func InitConfig(flags AppFlags) *Config {
	currDir := getCurrentDir()
	configPath := getConfigPath(currDir, flags)

	env, _ := loadEnvFile(currDir)
	yaml, _ := loadYmlFile(configPath)
	cfg := &Config{
		EnvParams:  env,
		FileParams: yaml,
		flags:      flags,
		Merged: MergedParams{
			LogInFile:  getLogInFile(),
			LogLevel:   getLogLevel(env, flags),
			LogDir:     getLogDir(env, flags),
			ConfigPath: configPath,
		},
	}

	return cfg
}

func getCurrentDir() string {
	var dir string
	currentPath, err := os.Getwd()
	if err != nil {
		panic("Cant read current path")
	}
	dir = filepath.Dir(currentPath)
	return dir
}

func loadEnvFile(dir string) (config EnvParams, err error) {
	viper.AddConfigPath(dir)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	readErr := viper.ReadInConfig()

	if err = viper.Unmarshal(&config); err != nil {
		// if key not present in .env file - then it will not be loaded from environment to struct
		panic("Fail mapping configuration parameters (env)")
	}
	if readErr != nil {
		// its': config == (EnvParams{})
		// try load all from natural environment
		// we don't want to use viper.Get("APP_ENV") because you need to remember all the variables
		// we need the config's IDE autocompleting
		if err := cleanenv.ReadEnv(&config); err != nil {
			panic("Fail mapping env params structure: " + err.Error())
		}
	} else {
		config.FileLoaded = true
	}

	return config, nil
}

func loadYmlFile(configPath string) (config FileParams, err error) {
	viper.SetConfigFile(configPath)

	viper.SetDefault("test_string", "qwerty")

	readErr := viper.ReadInConfig()

	if err = viper.Unmarshal(&config); err != nil {
		panic("Fail mapping configuration parameters (yaml)")
	}
	if readErr == nil {
		config.FileLoaded = true
	}

	return config, nil
}

func getConfigPath(currDir string, flags AppFlags) string {
	var configPath string
	if os.Getenv("APP_CONFIG_PATH") != "" {
		configPath = os.Getenv("APP_CONFIG_PATH")
	} else {
		if flags.ConfigPath.set == true {
			configPath = flags.ConfigPath.value
		} else {
			configPath = filepath.Join(currDir, DefaultConfigPath)
		}
	}
	return configPath
}

func getLogDir(env EnvParams, flags AppFlags) string {
	var logDir string
	if env.LogDir != "" {
		logDir = env.LogDir
	} else {
		if flags.LogDir.set == true {
			logDir = flags.LogDir.value
		} else {
			logDir = DefaultLogDir
		}
	}
	return logDir
}

func getLogInFile() bool {
	switch os.Getenv("LOG_IN_FILE_ENABLED") {
	case "":
		fallthrough
	case "false":
		fallthrough
	case "0":
		fallthrough
	case "no":
		return false
	case "1":
		fallthrough
	case "true":
		fallthrough
	case "yes":
		return true
	}
	return false
}

func getLogLevel(env EnvParams, flags AppFlags) string {
	var logLevel string
	if env.LogLevel != "" {
		logLevel = env.LogLevel
	} else {
		if flags.LogLevel.set == true {
			logLevel = flags.LogLevel.value
		} else {
			logLevel = DefaultLogLevel
		}
	}
	return logLevel
}
