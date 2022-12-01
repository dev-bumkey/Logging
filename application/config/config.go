package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/database/impl"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/kafka"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/server"
	"github.com/davecgh/go-spew/spew"
)

var Data *Config
var err error

type Config struct {
	DBConfig     *impl.Config
	ServerConfig *server.Config
	KafkaConfig  *kafka.Config

	AuthServerUrl string `env:"AUTH_SERVER_URL"` // for backward comportability

	LoggingLevel    string `env:"LOGGING_LEVEL" envDefault:"info"`
	LoggingFileUse  bool   `env:"LOGGING_FILE_USE" envDefault:"false"`
	LoggingFilePath string `env:"LOGGING_FILE_PATH" envDefault:"/var/log/cocktailcloud/alarm-collector.log"`

	// Debug Properties -- start
	RequestTrace bool `env:"REQUEST_TRACE" envDefault:"false"`
	DevMode      bool `env:"DEV_MODE" envDefault:"false"`
	Profiling    bool `env:"DEV_PROFILING" envDefault:"false"`
	SkipAuth     bool `env:"DEV_SKIP_AUTH" envDefault:"false"`
	SkipSave     bool `env:"DEV_SKIP_SAVE" envDefault:"false"`
	SkipSqlExec  bool `env:"DEV_SKIP_SQL_EXEC" envDefault:"false"`
	SkipQueue    bool `env:"DEV_SKIP_QUEUE" envDefault:"false"` // for backward comportability
	// Debug Properties -- end

	DoAlarmNotify           bool   `env:"DO_ALARM_NOTIFY" envDefault:"false"`
	AlarmNotificationApiUrl string `env:"ALARM_NOTIFICATION_API_URL"`

	IntervalForOrphan string `env:"INTERVAL_ORPHAN" envDefault:"3 day"`
	RetryLimit        int    `env:"RETRY_LIMIT"`
}

var Default = &Config{
	DBConfig: &impl.Config{
		Type:         "postgres",  //디비 엔진명
		Host:         "localhost", //디비 호스트
		Port:         "5432",      //디비 port
		DatabaseName: "",          //디비 네임스페이스명
		UserName:     "",          //디비 사용자
		Password:     "",          //디비 패스워드
		MaxIdleConns: 3,
		MaxOpenConns: 3,
		UseTls:       false,
	},
	ServerConfig: &server.Config{
		Host: "0.0.0.0",
		Port: "9311",
	},
	KafkaConfig: &kafka.Config{
		UseKafka:        false,
		BrokerAddresses: []string{"localhost:9092"},
		Topic:           []string{"alarms"},
		Group:           "alarm-consumers",
	},
	IntervalForOrphan: "3 day",
	RetryLimit:        5,
}

func init() {
	Data = Default
	if err = env.Parse(Data); err != nil {
		logger.Error("Setup file loading was not performed normally. \nSome features may work with default settings and return undesirable results.\n", err)
	}
	spew.Dump(*Data)
}

func Load() (*Config, error) {
	return Data, err
}
