package log

import (
	"fmt"
	"log"
	"sync"
)

var infoLogger, errorLogger *log.Logger

func Info(format string, valArrary ...interface{}) {
	_ = infoLogger.Output(2, fmt.Sprintf(format, valArrary...))
}

func Error(format string, valArrary ...interface{}) {
	_ = infoLogger.Output(2, fmt.Sprintf(format, valArrary...))
}

var writer *dailyFileWriter

func Config(outputFileName string) {
	writer = &dailyFileWriter{
		filename:    outputFileName,
		lastYearDay: -1,
		lock:        &sync.Mutex{},
	}
	infoLogger = log.New(writer, "[ INFO ]", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	errorLogger = log.New(writer, "[ ERROR ]", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}
