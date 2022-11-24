package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	cfg "github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/log/graylog"
	"github.com/spf13/viper"

	"github.com/kylelemons/godebug/pretty"
	"github.com/sirupsen/logrus"
)

type LoyaltyLogger struct {
	logrus *logrus.Logger

	path string
}

func NewLogger(path string) (logger LoyaltyLogger) {
	logger = LoyaltyLogger{
		logrus: logrus.New(),
		path:   path,
	}
	config := cfg.GetConfig()
	if config != nil {
		logger.LoadConfig(config)
	}
	return
}

func (l *LoyaltyLogger) LoadConfig(cfg *viper.Viper) {
	log_config := cfg.GetStringMapString("log")
	report := log_config["report"]
	if report == "true" {
		l.logrus.SetReportCaller(true)
	}
	format := log_config["format"]
	if format == "json" {
		l.logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			// PrettyPrint:     true,
		})
	} else {
		l.logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	level := log_config["level"]
	switch level {
	case "debug":
		l.logrus.SetLevel(logrus.DebugLevel)
		break
	case "info":
		l.logrus.SetLevel(logrus.InfoLevel)
		break
	case "trace":
		l.logrus.SetLevel(logrus.TraceLevel)
		break
	case "warn":
		l.logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		l.logrus.SetLevel(logrus.ErrorLevel)
		break
	case "panic":
		l.logrus.SetLevel(logrus.PanicLevel)
		break
	case "fatal":
		l.logrus.SetLevel(logrus.FatalLevel)
		break
	default:
		l.logrus.SetLevel(logrus.WarnLevel)
		break
	}
	logFile := log_config["logfile"]
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logrus.Fatal("error opening file: %v", err)
		}
		l.logrus.SetOutput(f)
	}

	//Config graylog
	graylog_server := log_config["graylog_server"]
	if graylog_server != "" {
		graylog_port := log_config["graylog_port"]
		if graylog_port == "" {
			graylog_port = "12201"
		}
		graylog_mode := log_config["graylog_mode"]
		if graylog_mode == "" {
			graylog_mode = "async"
		}
		graylog_url := fmt.Sprintf("%s:%s", graylog_server, graylog_port)
		var hook logrus.Hook
		if graylog_mode == "sync" {
			hook = graylog.NewGraylogHook(graylog_url, map[string]interface{}{})
		} else {
			hook = graylog.NewAsyncGraylogHook(graylog_url, map[string]interface{}{})
		}

		l.logrus.AddHook(hook)
	}
}

func (l *LoyaltyLogger) WithFields() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return l.logrus.WithField("file", filename).WithField("function", fn)
}

func (l *LoyaltyLogger) Println(message string, args ...interface{}) {
	l.WithFields().Infof(message, args...)
}

func (l *LoyaltyLogger) Printf(message string, args ...interface{}) {
	l.WithFields().Printf(message, args...)
}

func (l *LoyaltyLogger) Fatalf(message string, args ...interface{}) {
	l.WithFields().Fatalf(message, args...)
}

func (l *LoyaltyLogger) Infof(message string, args ...interface{}) {
	l.WithFields().Infof(message, args...)
}

func (l *LoyaltyLogger) Debugf(message string, args ...interface{}) {
	l.WithFields().Debugf(message, args...)
}

func (l *LoyaltyLogger) Warnf(message string, args ...interface{}) {
	l.WithFields().Warnf(message, args...)
}

func (l *LoyaltyLogger) Errorf(message string, args ...interface{}) {
	l.WithFields().Errorf(message, args...)
}

func (l *LoyaltyLogger) PrettyPrint(message string, args ...interface{}) {
	diffable := &pretty.Config{
		Diffable: true,
	}

	l.Infof(fmt.Sprintf("%s: %s\n", message, diffable.Sprint(args)))
}

func trace() (path string, functionName string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fmt.Println("frames")

	return fmt.Sprintf("%s:%d", frame.File, frame.Line), frame.Function
}
