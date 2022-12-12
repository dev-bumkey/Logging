// Package logger - Common Logging (inspired https://www.mountedthoughts.com/golang-logger-interface/)
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	newLog "log"
	"strings"

	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelInfo             = "info"
	InstanceZapLogger int = iota
)

// ===== [ Constants and Variables ] =====
type Logger interface {
	//Info(args ...interface{})
	Infof(format string, args ...interface{})
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

var (
	// A global variable so that log functions can be directly acccessed
	log                      Logger
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

// ===== [ Implementations ] =====

// Infof - Implements of Logger's Infof
func (zl *zapLogger) Infof(format string, args ...interface{}) {
	zl.sugaredLogger.Infof(format, args...)
}

type Config struct {
	// EnableConsole     bool
	// ConsoleJSONFormat bool
	// ConsoleLevel      string
	EnableFile     bool
	FileJSONFormat bool
	//FileLevel      string
	FileLocation string
}

func getEncoder(isJSON bool) zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func RegisterEncoder(isJSON bool) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:       "",
		LevelKey:      "",
		NameKey:       "",
		CallerKey:     "",
		MessageKey:    "M",
		StacktraceKey: "",
	}

	if isJSON {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// Info - Writes info level information.
func History(format string, args ...interface{}) {
	if log == nil {
		fmt.Printf(format, args...)
	} else {
		log.Infof(format, args...)

	}
}

// newZapLogger - returns an instance of Logger or Error for Zap
func NewZapLogger(alarmconf config.Config) (Logger, error) {
	cores := []zapcore.Core{}
	// if conf.EnableConsole {
	// 	level := getZapLevel(conf.ConsoleLevel)
	// 	writer := zapcore.Lock(os.Stdout)
	// 	core := zapcore.NewCore(getEncoder(conf.ConsoleJSONFormat), writer, level)
	// 	cores = append(cores, core)
	// }

	if alarmconf.EnableFile {
		level := zapcore.InfoLevel
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: alarmconf.FileLocation,
			MaxSize:  100,
			Compress: true,
			MaxAge:   28,
		})
		core := zapcore.NewCore(RegisterEncoder(false), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skip 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore, zap.AddCallerSkip(3), zap.AddCaller()).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
	// log = &zapLogger{
	// 	sugaredLogger: logger,
	// }
}

func AlarmHistoryFormat(alarmhistory *structure.AlarmHistory) {

	//alarmConfig := config.Config{}

	alarmJson := &alarmhistory
	result, err := json.Marshal(alarmJson)

	if err != nil {
		newLog.Fatalf("JSON marshaling failed: %s", err)
	}
	if config.Data.FileJSONFormat == "txt" {
		var txtFormat = structure.AlarmHistory{}
		err := json.Unmarshal([]byte(result), &txtFormat)
		if err != nil {
			fmt.Println("Failed to json.Unmarshal", err)
		}

		fmt.Printf("%s|%s|%s|%s|%s|%s|%s|%s", txtFormat.ClusterId, txtFormat.Alertname, txtFormat.RuleId, txtFormat.Severity, txtFormat.Status, txtFormat.StartsAt, txtFormat.EndsAt, strings.ReplaceAll(txtFormat.EngMsg, "\"", ""))
		txtResult := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", txtFormat.ClusterId, txtFormat.Alertname, txtFormat.RuleId, txtFormat.Severity, txtFormat.Status, txtFormat.StartsAt, txtFormat.EndsAt, strings.ReplaceAll(txtFormat.EngMsg, "\"", ""))
		History(txtResult)
	} else {
		fmt.Println(string(result))
		History(string(result))
	}

}

func NewTestLogger(conf config.Config, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := NewZapLogger(conf)

		if err != nil {
			return err
		}
		log = logger
		return nil
	default:
		return errInvalidLoggerInstance
	}
}

func NewAlarmLogger(conf *config.Config) error {
	return NewTestLogger(config.Config{
		// ConsoleJSONFormat: true,
		EnableFile:     conf.EnableFile,
		FileJSONFormat: conf.FileJSONFormat,
		FileLocation:   conf.FileLocation,
	}, InstanceZapLogger)
}
