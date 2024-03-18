package logger

import "log"

type TodoLogger struct {
	infoLogger    log.Logger
	warningLogger log.Logger
	errorLogger   log.Logger
}

func (logger TodoLogger) Info() {

}
