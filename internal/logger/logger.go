package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}
