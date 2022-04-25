package shared

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type LoggerCfg struct {
	Level  log.Level `json:"level"`
	File   LogFile   `json:"file"`
	SysLog SysLog    `json:"syslog"`
}

type LogFile struct {
	Enabled    bool   `json:"enabled"`
	FileName   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}

type SysLog struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"address"`
	Network string `json:"network"`
	Tag     string `json:"tag"`
}

var logger *log.Logger

func InitLogger(conf LoggerCfg) {
	logger = &log.Logger{
		Out:          os.Stdout,
		Hooks:        nil,
		Formatter:    &log.JSONFormatter{},
		ReportCaller: false,
		Level:        conf.Level,
		ExitFunc:     os.Exit,
	}
	if conf.File.Enabled {
		if len(conf.File.FileName) == 0 {
			panic(fmt.Errorf("parameter FileName not defined"))
		}
		rotateHook, err := NewRotateFileHook(RotateFileConfig{
			Filename:   conf.File.FileName,
			MaxSize:    conf.File.MaxSize,
			MaxBackups: conf.File.MaxBackups,
			MaxAge:     conf.File.MaxAge,
			Level:      conf.Level,
			Formatter:  &log.JSONFormatter{},
		})
		if err != nil {
			logger.Info("Failed to log to file")
			return
		}
		logger.Hooks = log.LevelHooks{}
		logger.AddHook(rotateHook)
	}
	if conf.SysLog.Enabled {
		syslog, err := NewSyslogHook(conf)
		if err != nil {
			logger.Info("Failed to log to syslog")
			return
		}
		if logger.Hooks == nil {
			logger.Hooks = log.LevelHooks{}
		}
		logger.AddHook(syslog)
	}
}

func GetLogger() *log.Logger {
	return logger
}
