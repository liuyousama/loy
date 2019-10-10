package log

import (
	"fmt"
	"github.com/liuyousama/loy/log/text_generator"
	"github.com/liuyousama/loy/log/text_handler"
)

var l = &logger{}

type logger struct {
	level     int
	generator text_generator.TextGenerator
	handlers  []text_handler.Handler
}

func Fatal(text string) {
	l.fatal(text)
}
func (log *logger) fatal(text string) {
	text = log.generator.Generate("", text, fatalLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, fatalLevel, log.level)
	}
}

func Fatalf(format string, a ...interface{})  {
	l.fatalf(format, a...)
}
func (log *logger) fatalf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate("", text, fatalLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, fatalLevel, log.level)
	}
}

func Error(text string)  {
	l.error(text)
}
func (log *logger) error(text string) {
	text = log.generator.Generate("", text, errorLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, errorLevel, log.level)
	}
}

func Errorf(format string, a ...interface{})  {
	l.errorf(format, a...)
}
func (log *logger) errorf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate("", text, errorLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, errorLevel, log.level)
	}
}

func Debug(text string)  {
	l.debug(text)
}
func (log *logger) debug(text string) {
	text = log.generator.Generate("", text, debugLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, debugLevel, log.level)
	}
}

func Debugf(format string, a ...interface{})  {
	l.debugf(format, a...)
}
func (log *logger) debugf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate("", text, debugLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, debugLevel, log.level)
	}
}

func Info(text string)  {
	l.info(text)
}
func (log *logger) info(text string) {
	text = log.generator.Generate("", text, infoLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, infoLevel, log.level)
	}
}

func Infof(format string, a ...interface{})  {
	l.infof(format, a...)
}
func (log *logger) infof(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	text = log.generator.Generate("", text, infoLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, infoLevel, log.level)
	}
}

func InfoWithTag(tag, text string) {
	l.infoWithTag(tag, text)
}
func (log *logger) infoWithTag(tag, text string) {
	text = log.generator.Generate(tag, text, infoLevelText)
	for _, handler := range log.handlers {
		handler.HandleText(text, infoLevel, log.level)
	}
}