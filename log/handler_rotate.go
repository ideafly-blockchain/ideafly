package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateConfig struct {
	LogDir     string `json:"log_dir"`
	Filename   string `json:"filename"` // file name
	MaxAge     int    `json:"max_age"`  // max age
	MaxSize    int    `json:"max_size"` // MB
	MaxBackups int    `json:"max_backups"`
}

var defaultConfig = &RotateConfig{
	LogDir:     "logs",
	Filename:   "chain.log",
	MaxSize:    100,
	MaxAge:     7,
	MaxBackups: 10,
}

func NewRotateConfig() *RotateConfig {
	conf := *defaultConfig
	return &conf
}

func NewFileRotateHandler(config *RotateConfig, format Format) Handler {
	if err := config.setup(); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	logDir := config.LogDir
	if !filepath.IsAbs(logDir) {
		logDir, _ = filepath.Abs(logDir)
	}
	log := lumberjack.Logger{
		Filename:   path.Join(logDir, config.Filename),
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		LocalTime:  true,
		Compress:   true, // disabled by default
	}

	h := StreamHandler(&log, format)

	return FuncHandler(func(r *Record) error {
		return h.Log(r)
	})
}

func (c *RotateConfig) setup() error {
	if len(c.LogDir) == 0 {
		panic("Failed to parse logger folder:" + c.LogDir + ".")
	}

	if err := os.MkdirAll(c.LogDir, 0700); err != nil {
		panic("Failed to create logger folder:" + c.LogDir + ". err:" + err.Error())
	}
	return nil
}
