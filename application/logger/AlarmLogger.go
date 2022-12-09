// Package logger - Common Logging (inspired https://www.mountedthoughts.com/golang-logger-interface/)
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	newLog "log"

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
	Info(args ...interface{})
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

// Info - Implements of Logger's Info
func (zl *zapLogger) Info(args ...interface{}) {
	zl.sugaredLogger.Info(args...)
}

type Config struct {
	// EnableConsole     bool
	// ConsoleJSONFormat bool
	// ConsoleLevel      string
	EnableFile     bool
	FileJSONFormat bool
	FileLevel      string
	FileLocation   string
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

// Info - Writes info level information.
func History(args ...interface{}) {
	if log == nil {
		fmt.Println(args...)
	} else {
		log.Info(args)

	}
}

// newZapLogger - returns an instance of Logger or Error for Zap
func NewZapLogger(enable bool, path string) {
	cores := []zapcore.Core{}
	// if conf.EnableConsole {
	// 	level := getZapLevel(conf.ConsoleLevel)
	// 	writer := zapcore.Lock(os.Stdout)
	// 	core := zapcore.NewCore(getEncoder(conf.ConsoleJSONFormat), writer, level)
	// 	cores = append(cores, core)
	// }

	if enable {
		level := zapcore.InfoLevel
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: path,
			MaxSize:  100,
			Compress: true,
			MaxAge:   28,
		})
		core := zapcore.NewCore(getEncoder(false), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skip 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore, zap.AddCallerSkip(3), zap.AddCaller()).Sugar()

	// return &zapLogger{
	// 	sugaredLogger: logger,
	// }, nil
	log = &zapLogger{
		sugaredLogger: logger,
	}
}

func AlarmHistoryFormat(format string, alarmhistory *structure.AlarmHistory) {

	alarmJson := &alarmhistory
	result, err := json.Marshal(alarmJson)

	if err != nil {
		newLog.Fatalf("JSON marshaling failed: %s", err)
	}
	if format == "json" {
		fmt.Println(string(result))
		History(string(result))
	} else {
		var txtFormat = structure.AlarmHistory{}
		err := json.Unmarshal([]byte(result), &txtFormat)
		if err != nil {
			fmt.Println("Failed to json.Unmarshal", err)
		}

		// fmt.Println(txtFormat)
		fmt.Printf("%s|%s|%s|%s|%s|%s|%s|%s", txtFormat.ClusterId, txtFormat.Alertname, txtFormat.RuleId, txtFormat.Severity, txtFormat.Status, txtFormat.StartsAt, txtFormat.EndsAt, txtFormat.Description)
		txtResult := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", txtFormat.ClusterId, txtFormat.Alertname, txtFormat.RuleId, txtFormat.Severity, txtFormat.Status, txtFormat.StartsAt, txtFormat.EndsAt, txtFormat.Description)
		History(txtResult)
	}

}
