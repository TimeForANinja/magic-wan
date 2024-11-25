package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

const logfile = "/var/log/magic-wan.log"

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
// Provided as solution here: https://github.com/sirupsen/logrus/issues/678#issuecomment-362569561
type WriterHook struct {
	Writer    io.Writer
	LogLevels []log.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []log.Level {
	return hook.LogLevels
}

func configureLogging() *os.File {
	// Configure logging to file
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	panicOn(err)

	// prep global logger
	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// add custom loggers to have two destinations with different log level
	log.AddHook(&WriterHook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
		},
	})
	log.AddHook(&WriterHook{
		Writer: file,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
			log.DebugLevel,
		},
	})

	return file
}
