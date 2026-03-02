/*
Package config конфигурация

	Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
	Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например, http://localhost:8000/qsd54gFg).
	Флаг -l отвечает за уровень логирования (значение по умолчанию: "Info")
	Флаг -f путь до файла, куда сохраняются данные в формате JSON (значение по умолчанию "./storage.json")
*/
package config

import (
	"encoding/json"
	"flag"
	"os"
	"sync"
)

// ShortLen длина сокращенного кода.
const ShortLen = 6

type sourceType int

const (
	_ sourceType = iota
	sourceEnv
	sourceFlag
	sourceConfig
	sourceDefault
)

type valueTypeType int

const (
	_ valueTypeType = iota
	valueString
	valueBool
)

type ConfigValue struct {
	defaultValue any
	flagName     []string
	envName      string
	description  string

	value     any
	valueType valueTypeType

	source sourceType
}

type configType string

const (
	configServerAddress = "server_address"
	configBaseAddress   = "base_url"
	configLogLevel      = "log_level"
	configFileName      = "file_storage_path"
	configConnAddress   = "database_dsn"
	configAuditFile     = "audit_file"
	configAuditAddress  = "audit_url"
	configEnableHttps   = "enable_https"
	configConfigFile    = "config_file"
)

/*
	type jsonConfig struct {
	    serverAddress string `json:"server_address"` // : "localhost:8080", // аналог переменной окружения SERVER_ADDRESS или флага -a
	    baseAddress   string `json:"base_url"` // : "http://localhost", // аналог переменной окружения BASE_URL или флага -b
	    fileName      string `json:"file_storage_path"` // : "/path/to/file.db", // аналог переменной окружения FILE_STORAGE_PATH или флага -f
	    connAddress   string `json:"database_dsn` // ": "", // аналог переменной окружения DATABASE_DSN или флага -d
	    enableHTTPS   bool `json:"enable_https"` // : true, // аналог переменной окружения ENABLE_HTTPS или флага -s
	    logLevel      string `json:"log_level"` // : "Info", // аналог переменной окружения LOG_LEVEL или флага -l
	    auditFile     string `json:"audit_file"` // : "", // аналог переменной окружения AUDIT_FILE или флага -audit-file
	    auditAddress  string `json:"audit_url"` // : "" // аналог переменной окружения AUDIT_URL или флага -audit-url
	}
*/
const (
	defaultServerAddress = ":8080"
	defaultBaseAddress   = "http://localhost:8080"
	defaultLogLevel      = "Info"
)

var (
	configData = map[configType]*ConfigValue{
		configServerAddress: {defaultValue: defaultServerAddress, flagName: []string{"a"}, envName: "SERVER_ADDRESS", description: "", valueType: valueString},
		configBaseAddress:   {defaultValue: defaultBaseAddress, flagName: []string{"b"}, envName: "BASE_URL", description: "", source: sourceDefault, valueType: valueString},
		configLogLevel:      {defaultValue: defaultLogLevel, flagName: []string{"l"}, envName: "LOG_LEVEL", description: "", valueType: valueString},
		configFileName:      {defaultValue: "", flagName: []string{"f"}, envName: "FILE_STORAGE_PATH", description: "", valueType: valueString},
		configConnAddress:   {defaultValue: "", flagName: []string{"d"}, envName: "DATABASE_DSN", description: "", valueType: valueString},
		configAuditFile:     {defaultValue: "", flagName: []string{"audit-file"}, envName: "AUDIT_FILE", description: "", valueType: valueString},
		configAuditAddress:  {defaultValue: "", flagName: []string{"audit-url"}, envName: "AUDIT_URL", description: "", valueType: valueString},
		configEnableHttps:   {defaultValue: false, flagName: []string{"s"}, envName: "ENABLE_HTTPS", description: "", valueType: valueBool},
		configConfigFile:    {defaultValue: "", flagName: []string{"c", "config"}, envName: "CONFIG", description: "файл конфигурации", valueType: valueString},
	}
)

/*
const (
	defaultServerAddress = ":8080"
	flagServerAddress    = "a"
	envServerAddress     = "SERVER_ADDRESS"
	descServerAddress    = "адрес запуска HTTP-сервера"

	defaultBaseAddress = "http://localhost:8080"
	flagBaseAddress    = "b"
	envBaseAddress     = "BASE_URL"
	descBaseAddress    = "базовый адрес результирующего сокращённого URL"

	defaultLogLevel = "Info"
	flagLogLevel    = "l"
	envLogLevel     = "LOG_LEVEL"
	descLogLevel    = "уровень логирования"

	defaultFileName = ""
	flagFileName    = "f"
	envFileName     = "FILE_STORAGE_PATH"
	descFileName    = "файл для хранения сокращенных адресов"

	defaultConnAddress = ""
	flagConnAddress    = "d"
	envConnAddress     = "DATABASE_DSN"
	descConnAddress    = "строка с адресом подключения к БД"

	defaultAuditFile = ""
	flagAuditFile    = "audit-file"
	envAuditFile     = "AUDIT_FILE"
	descAuditFile    = "путь к файлу-приёмнику, в который сохраняются логи аудита"

	defaultAuditAddress = ""
	flagAuditAddress    = "audit-url"
	envAuditAddress     = "AUDIT_URL"
	descAuditAddress    = "полный URL удаленного сервера-приёмника, куда отправляются логи аудита"

	defaultEnableHTTPS = false
	flagEnableHTTPS    = "s"
	envEnableHTTPS     = "ENABLE_HTTPS"
	descEnableHTTPS    = "использовать HTTPS"
)
*/

// DefaultConfig дефолтовый конфиг для тестов.
/*
var DefaultConfig = &Config{
	serverAddress: defaultServerAddress,
	baseAddress:   defaultBaseAddress,
	logLevel:      defaultLogLevel,
	fileName:      defaultFileName,
	connAddress:   defaultConnAddress,
	auditFile:     defaultAuditFile,
	auditAddress:  defaultAuditAddress,
}
*/

// Config структура для хранения конфига.
type Config struct {
	list map[configType]*ConfigValue
	/*
		serverAddress string
		baseAddress   string
		logLevel      string
		fileName      string
		connAddress   string
		auditFile     string
		auditAddress  string
		enableHTTPS   bool
	*/
}

/*
type ConfigFile struct {
	serverAddress string
	baseAddress   string
		logLevel      string
		fileName      string
		connAddress   string
		auditFile     string
		auditAddress  string
		enableHTTPS   bool

}
*/

/*
var (
	config = map[configType]*ConfigValue{
		configServerAddress: {defaultValue: ":8080", flagName: "a", envName: "SERVER_ADDRESS", description: "", valueType: valueString},
		configBaseAddress:   {defaultValue: "http://localhost:8080", flagName: "b", envName: "BASE_URL", description: "", source: sourceDefault, valueType: valueString},
		configLogLevel:      {defaultValue: "Info", flagName: "l", envName: "LOG_LEVEL", description: "", valueType: valueString},
		configFileName:      {defaultValue: "", flagName: "f", envName: "FILE_STORAGE_PATH", description: "", valueType: valueString},
		configConnAddress:   {defaultValue: "", flagName: "d", envName: "DATABASE_DSN", description: "", valueType: valueString},
		configAuditFile:     {defaultValue: "", flagName: "audit-file", envName: "AUDIT_FILE", description: "", valueType: valueString},
		configAuditAddress:  {defaultValue: "", flagName: "audit-url", envName: "AUDIT_URL", description: "", valueType: valueString},
		configEnableHttps:   {defaultValue: "", flagName: "s", envName: "ENABLE_HTTPS", description: "использовать HTTPS", valueType: valueBool},
		configConfigFile:    {defaultValue: "", flagName: "c", envName: "CONFIG", description: "файл конфигурации", valueType: valueString},
	}
)
*/

var fn sync.Once

// New создает и инициализирует структуру с конфигурацией.
func New() *Config {

	fn.Do(func() {

		for _, v := range configData {
			v.value = v.defaultValue
			v.source = sourceDefault
			switch v.valueType {

			case valueString:
				for _, f := range v.flagName {
					_ = flag.String(f, v.defaultValue.(string), v.description)
				}
			case valueBool:
				for _, f := range v.flagName {
					_ = flag.Bool(f, v.defaultValue.(bool), v.description)
				}
			}
			if respEnv, ok := os.LookupEnv(v.envName); ok && respEnv != "" {
				v.value = respEnv
				v.source = sourceEnv
			}
		}

		flag.Parse()

		flag.Visit(func(flagValue *flag.Flag) {
			if v, ok := configData[configType(flagValue.Name)]; ok {
				v.source = sourceFlag
				switch v.valueType {
				case valueString:
					v.value = flagValue.Value.String()
				case valueBool:
					v.value = flagValue.Value.String() == "true"
				}
			}
		})

		if confFile, ok := configData[configConfigFile]; ok {
			if data, err := os.ReadFile(confFile.value.(string)); err == nil {
				val := make(map[string]any)
				if err = json.Unmarshal(data, &val); err == nil {
					for k, v := range val {
						if item, ok := configData[configType(k)]; ok {
							if item.source == sourceDefault {
								item.source = sourceConfig
								switch item.valueType {
								case valueString:
									item.value = v.(string)
								case valueBool:
									item.value = v.(bool)
								}
							}
						}
					}
				}
			}
		}

		/*
			serverAddress := setParamString(envServerAddress, flagServerAddress, defaultServerAddress, descServerAddress)
			baseAddress := setParamString(envBaseAddress, flagBaseAddress, defaultBaseAddress, descBaseAddress)
			logLevel := setParamString(envLogLevel, flagLogLevel, defaultLogLevel, descLogLevel)
			fileName := setParamString(envFileName, flagFileName, defaultFileName, descFileName)
			connAddress := setParamString(envConnAddress, flagConnAddress, defaultConnAddress, descConnAddress)
			auditFile := setParamString(envAuditFile, flagAuditFile, defaultAuditFile, descAuditFile)
			auditAddress := setParamString(envAuditAddress, flagAuditAddress, defaultAuditAddress, descAuditAddress)
			enableHTTPS := setParamBool(envEnableHTTPS, flagEnableHTTPS, defaultEnableHTTPS, descEnableHTTPS)

			flag.Parse()

			x = &Config{
				serverAddress: validateServerAddress(*serverAddress, defaultServerAddress),
				baseAddress:   validateBaseAddress(*baseAddress, defaultBaseAddress),
				logLevel:      *logLevel,
				fileName:      *fileName,
				connAddress:   *connAddress,
				auditFile:     *auditFile,
				auditAddress:  validateBaseAddress(*auditAddress, defaultAuditAddress),
				enableHTTPS:   *enableHTTPS,
			}
		*/
	})
	return &Config{list: configData}
}

/*

 */

// setParamString reurn usage, value
/*
func setParam(env, flagName, defaultAddress, description string) (usage bool, respFlag *flag) {
	f := flag.Lookup(flagName)
	if f != nil {
		respFlag = &f.Value
	}
	respFlag := flag.String(flagName, defaultAddress, description)
	if respEnv, ok := os.LookupEnv(env); ok && respEnv != "" {
		return sourceEnv, &respEnv
	}
	return respFlag
}


func setParamBool(env, flagName string, defaultParam bool, description string) *bool {
	respFflag := flag.Bool(flagName, defaultParam, description)
	if x, ok := os.LookupEnv(env); ok && x != "" {
		respEnv := true
		return &respEnv
	}
	return respFflag
}
*/
/*
func validateServerAddress(address, defaultAddress string) string {
	addrList := strings.Split(address, ":")
	if len(addrList) < 1 || len(addrList) > 2 || len(addrList) == 1 && addrList[0] == "" {
		return defaultAddress
	}

	if len(addrList) < 2 || addrList[1] == "" {
		return addrList[0]
	}
	return net.JoinHostPort(addrList[0], addrList[1])
}

func validateBaseAddress(address, defaultAddress string) string {
	u, err := url.Parse(address)
	if err != nil {
		return defaultAddress
	}

	if u.Scheme == "" || u.Host == "" {
		return defaultAddress
	}

	return u.String()
}
*/

func (c *Config) getString(name configType) string {
	if v, ok := c.list[name]; ok && v != nil && v.valueType == valueString && v.value != nil {
		return v.value.(string)
	}
	return ""
}

func (c *Config) getBool(name configType) bool {
	if v, ok := c.list[name]; ok && v != nil && v.valueType == valueBool && v.value != nil {
		return v.value.(bool)
	}
	return false

}

// GetBaseAddress за адрес запуска HTTP-сервера.
func (c *Config) GetBaseAddress() string {
	return c.getString(configBaseAddress)
}

// GetServerAddress базовый адрес результирующего сокращённого URL.
func (c *Config) GetServerAddress() string {
	return c.getString(configServerAddress)
}

// GetLogLevel уровень логирования.
func (c *Config) GetLogLevel() string {
	return c.getString(configLogLevel)
}

// GetFileName путь до файла, куда сохраняются данные в формате JSON.
func (c *Config) GetFileName() string {
	return c.getString(configFileName)
}

// GetConnAddress адрес подключения к БД
func (c *Config) GetConnAddress() string {
	return c.getString(configConnAddress)
}

// GetAuditFile путь к файлу-приёмнику, в который сохраняются логи аудита.
func (c *Config) GetAuditFile() string {
	return c.getString(configAuditFile)
}

// GetAuditAddress полный URL удаленного сервера-приёмника, куда отправляются логи аудита.
func (c *Config) GetAuditAddress() string {
	return c.getString(configAuditAddress)
}

// IsEnableHTTPS использовать HTTPS
func (c *Config) IsEnableHTTPS() bool {
	return c.getBool(configEnableHttps)
}
