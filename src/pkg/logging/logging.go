package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	*logrus.Entry
}

var logger *logrus.Entry

func Get() *Logger {
	return &Logger{logger}
}

func init() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	logger = logrus.NewEntry(log)
}
