package logging_linux

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

func New(networks, server, processTag, action, logFileName string) *SyslogLinux {
	logCh := make(chan LogEntry)

	return &SyslogLinux{
		Networks:    networks,
		Server:      server,
		ProcessTag:  processTag,
		Action:      action,
		LogFileName: logFileName,
		LogChannel:  logCh,
	}
}

func writeToFile(logger *SyslogLinux, severity, msg string) {
	if logger.LogFileName == "" {
		logger.LogFileName = fmt.Sprintf("generic_%s.log", logger.ProcessTag)
	}

	f, err := os.OpenFile(logger.LogFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	log.SetOutput(f)
	log.Println(severity + ":" + msg)
}

func writeToSyslog(logger *SyslogLinux, severity, msg string) {
	sysLog, err := syslog.Dial(logger.Networks, logger.Server, syslog.LOG_INFO|syslog.LOG_DAEMON, logger.ProcessTag)
	if err != nil {
		log.Fatal(err)
	}

	if severity == "info" {
		if err = sysLog.Info(msg); err != nil {
			log.Fatal(err)
		}
	} else if severity == "warning" {
		if err = sysLog.Warning(msg); err != nil {
			log.Fatal(err)
		}
	} else if severity == "error" {
		if err = sysLog.Err(msg); err != nil {
			log.Fatal(err)
		}
	} else {
		if err = sysLog.Info(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func (logger *SyslogLinux) Log(logLevel, msg string) {
	logger.LogChannel <- LogEntry{
		LogLevel: logLevel,
		Msg:      msg,
	}
}

func (logger *SyslogLinux) WriteLog() {
	for msg := range logger.LogChannel {
		switch logger.Action {
		case "syslog":
			writeToSyslog(logger, msg.LogLevel, msg.Msg)
		case "file":
			writeToFile(logger, msg.LogLevel, msg.Msg)
		case "all":
			writeToSyslog(logger, msg.LogLevel, msg.Msg)
			writeToFile(logger, msg.LogLevel, msg.Msg)
		default:
			writeToSyslog(logger, "info", msg.Msg)
		}
	}
}
