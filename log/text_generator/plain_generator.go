package text_generator

import (
	"fmt"
	"time"
)

type PlainGenerator struct {

}

func init()  {
	TextGenerators["plain"] = new(PlainGenerator)
}

func (*PlainGenerator)Generate(tag string, content string, level string) string  {
	fun, file, line := GetCaller()

	str := fmt.Sprintf("%s -- %s:%s:%d 【%s】%s-%s",
		time.Now().Format("2006-01-02 15:04:05"),
		file, fun, line,
		level, tag, content)

	return str
}
