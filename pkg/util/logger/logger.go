/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger handles all logger from application
	Logger struct {
		*zap.SugaredLogger
	}

	// LoggingFn is generic logging function with some additonal context
	LoggingFn func(level logLevel, msg string, keysAndValues ...interface{})

	logLevel string
)

var DefaultLogger *Logger
var debug bool

const (
	DebugLevel logLevel = "DEBUG"
	InfoLevel  logLevel = "INFO"
	WarnLevel  logLevel = "WARN"
	ErrorLevel logLevel = "ERROR"
)

/*
func init() {
	configChan := make(chan struct{})
	setting.RegNotifyChannel(configChan)

	go func() {
		for {
			select {
			case _, ok := <-configChan:
				{
					if !ok {
						return
					} else {
						Build()
					}
				}
			}
		}
	}()
}
*/

func init() {
	setting.RegSyncNotify(func() error {
		return Build()
	})
}

func Build() error {
	conf := setting.Config
	debug = conf.Debug
	var config *zap.Config
	if !conf.Debug {
		znpc := zap.NewProductionConfig()
		config = &znpc

		//ProductionConfig put log file in log.logpath
		config.OutputPaths = []string{conf.Log.LogPath}
		config.ErrorOutputPaths = []string{conf.Log.LogPath}
	} else {
		zndc := zap.NewDevelopmentConfig()
		config = &zndc

		if conf.Log.Logformat == "json" {
			config.Encoding = "json"
		} else {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	}

	//config.DisableStacktrace = true
	//config.Development = false
	config.DisableCaller = true

	logger, err := config.Build()
	if err != nil {
		return err
	}
	defer logger.Sync()
	DefaultLogger = &Logger{logger.Sugar()}
	return nil
}

/*
ContextLoggingFn creates a LoggingFn to be used in
places that do not necessarily need access to the gin context
*/
func (logger *Logger) ContextLoggingFn(c *gin.Context) LoggingFn {
	return func(level logLevel, msg string, keysAndValues ...interface{}) {
		switch level {
		case DebugLevel:
			logger.Debugc(c, msg, keysAndValues...)
		case InfoLevel:
			logger.Infoc(c, msg, keysAndValues...)
		case WarnLevel:
			logger.Warnc(c, msg, keysAndValues...)
		case ErrorLevel:
			logger.Errorc(c, msg, keysAndValues...)
		}
	}
}

func (logger *Logger) Print(v ...interface{}) {
	format_msg := gorm.LogFormatter(v...)
	msg := fmt.Sprintf("%v", format_msg)

	if debug {
		logger.Infof(msg)
	} else {
		logger.Errorf(msg)
	}
}

// Debugc wraps Debugw provided by zap, adding data from gin request context
func (logger *Logger) Debugc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Debugw(msg, keysAndValues...)
}

// Infoc wraps Infow provided by zap, adding data from gin request context
func (logger *Logger) Infoc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Infow(msg, keysAndValues...)
}

// Warnc wraps Warnw provided by zap, adding data from gin request context
func (logger *Logger) Warnc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Warnw(msg, keysAndValues...)
}

// Errorc wraps Errorw provided by zap, adding data from gin request context
func (logger *Logger) Errorc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Errorw(msg, keysAndValues...)
}

// transformLogcArgs prefixes msg with RequestCount and adds RequestId to keysAndValues
func transformLogcArgs(c *gin.Context, msg string, keysAndValues []interface{}) (string, []interface{}) {
	if reqCount, exists := c.Get("requestcount"); exists {
		msg = fmt.Sprintf("[%s] %s", reqCount, msg)
		keysAndValues = append(keysAndValues, "reqID", c.MustGet("requestid"))
	}
	return msg, keysAndValues
}

func init() {
	logrus.SetLevel(logrus.WarnLevel) // silence logs from zsais/go-gin-prometheus
}

/*
func init() {
	Log = logrus.New()
	Log.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: true}
}
*/
