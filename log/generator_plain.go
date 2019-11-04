package log

import (
	"fmt"
	"strings"
	"time"
)

type PlainGenerator struct {

}

//type check
var _ TextGenerator = &PlainGenerator{}

func (*PlainGenerator)generate(content string, level Level) string  {
	fun, file, line := getCaller()

	str := fmt.Sprintf("%s --【%s】%s %s:%s:%d ",
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ToUpper(level.GetLevelText()),  content,
		file, fun, line, )

	return str
}

func NewPlainGenerator() *PlainGenerator {
	return new(PlainGenerator)
}