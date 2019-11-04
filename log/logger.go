package log

import (
	"fmt"
)

var l = &logger{
	level:DEBUG,
}

type logger struct {
	level     Level
	generator TextGenerator
	handlers  []Handler
}

func LoadHandler(handlers ...Handler) (err error) {
	return l.loadHandler(handlers...)
}
func (log *logger)loadHandler(handlers ...Handler) (err error) {
	for _, handler := range handlers {
		err = handler.load()
		if err != nil {
			return err
		}
		log.handlers = append(log.handlers, handler)
	}

	return
}

func LoadGenerator(generator TextGenerator) {
	l.loadGenerator(generator)
}
func (log *logger)loadGenerator(generator TextGenerator) {
	log.generator = generator
}

func SetLevel(level Level) {
	l.setLevel(level)
}
func (log *logger)setLevel(level Level) {
	if level!=INFO&&level!=DEBUG&&level!=ERROR&&level!=FATAL {
		level = DEBUG
	}

	log.level = level
}

func Fatal(text string) {
	l.fatal(text)
}
func (log *logger) fatal(text string) {
	if log.level > FATAL {
		return
	}
	text = log.generator.generate(text, FATAL)
	for _, handler := range log.handlers {
		handler.handleText(text)
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
	text = log.generator.generate(text, FATAL)
	for _, handler := range log.handlers {
		handler.handleText(text)
	}
}

func Error(text string)  {
	l.error(text)
}
func (log *logger) error(text string) {
	if log.level > ERROR {
		return
	}
	text = log.generator.generate(text, ERROR)
	for _, handler := range log.handlers {
		handler.handleText(text)
	}
}

func Errorf(format string, a ...interface{})  {
	l.errorf(format, a...)
}
func (log *logger) errorf(format string, a ...interface{}) {
	if log.level > ERROR {
		return
	}
	text := fmt.Sprintf(format, a...)
	text = log.generator.generate(text, ERROR)
	for _, handler := range log.handlers {
		handler.handleText(text)
	}
}

func Debug(text string)  {
	l.debug(text)
}
func (log *logger) debug(text string) {
	if log.level > DEBUG {
		return
	}
	text = log.generator.generate(text, DEBUG)
	for _, handler := range log.handlers {
		handler.handleText(text)
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
	text = log.generator.generate(text, DEBUG)
	for _, handler := range log.handlers {
		handler.handleText(text)
	}
}

func Info(text string)  {
	l.info(text)
}
func (log *logger) info(text string) {
	if log.level > INFO {
		return
	}
	text = log.generator.generate(text, INFO)
	for _, handler := range log.handlers {
		handler.handleText(text)
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
	text = log.generator.generate(text, INFO)
	for _, handler := range log.handlers {
		handler.handleText(text)
	}
}