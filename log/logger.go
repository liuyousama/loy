package log

import (
	"fmt"
)

var l = &logger{}

type logger struct {
	level     Level
	generator TextGenerator
	handlers  []Handler
}

func (log *logger)LoadHandler(handlers ...Handler) (err error) {
	for _, handler := range handlers {
		err = handler.Load()
		if err != nil {
			return err
		}
		log.handlers = append(log.handlers, handler)
	}

	return
}

func (log *logger)LoadGenerator(generator TextGenerator) {
	log.generator = generator
}

func Fatal(text string) {
	l.fatal(text)
}
func (log *logger) fatal(text string) {
	if log.level > FATAL {
		return
	}
	text = log.generator.Generate(text, FATAL)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Fatalf(format string, a ...interface{})  {
	l.fatalf(format, a...)
}
func (log *logger) fatalf(format string, a ...interface{}) {
	if log.level > FATAL {
		return
	}
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate(text, FATAL)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Error(text string)  {
	l.error(text)
}
func (log *logger) error(text string) {
	if log.level > DEBUG {
		return
	}
	text = log.generator.Generate(text, ERROR)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Errorf(format string, a ...interface{})  {
	l.errorf(format, a...)
}
func (log *logger) errorf(format string, a ...interface{}) {
	if log.level > DEBUG {
		return
	}
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate(text, ERROR)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Debug(text string)  {
	l.debug(text)
}
func (log *logger) debug(text string) {
	if log.level > DEBUG {
		return
	}
	text = log.generator.Generate(text, DEBUG)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Debugf(format string, a ...interface{})  {
	l.debugf(format, a...)
}
func (log *logger) debugf(format string, a ...interface{}) {
	if log.level > DEBUG {
		return
	}
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate(text, DEBUG)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Info(text string)  {
	l.info(text)
}
func (log *logger) info(text string) {
	if log.level > INFO {
		return
	}
	text = log.generator.Generate(text, INFO)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}

func Infof(format string, a ...interface{})  {
	l.infof(format, a...)
}
func (log *logger) infof(format string, a ...interface{}) {
	if log.level > INFO {
		return
	}
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate(text, INFO)
	for _, handler := range log.handlers {
		handler.HandleText(text)
	}
}