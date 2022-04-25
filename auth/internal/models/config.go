package models

import (
	"github.com/Ypxd/diplom/shared"
	"time"
)

type Config struct {
	Redis  Redis            `json:"redis"`
	Server Server           `json:"server"`
	Logger shared.LoggerCfg `json:"log_params"`
	DB     DB               `json:"db"`
}

type Redis struct {
	Address  string        `json:"address,omitempty"`
	Password string        `json:"password,omitempty"`
	RTimeout time.Duration `json:"r_timeout,omitempty"`
}

type Server struct {
	Host     string        `json:"host,omitempty"`
	Port     uint          `json:"port,omitempty"`
	RTimeout time.Duration `json:"r_timeout,omitempty"`
	WTimeout time.Duration `json:"w_timeout,omitempty"`
}

type Logger struct {
	Level  string `json:"level"`
	File   File   `json:"file"`
	Syslog Syslog `json:"syslog"`
}
type File struct {
	Enabled    bool   `json:"enabled"`
	FileName   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}

type Syslog struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"address"`
	Network string `json:"network"`
	Tag     string `json:"tag"`
}

type DB struct {
	Address    string `json:"address"`
	DBName     string `json:"db_name"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Trace      bool   `json:"trace"`
	DriverName string `json:"driver_name"`
}
