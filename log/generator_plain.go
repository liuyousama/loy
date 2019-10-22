package log

import (
	"fmt"
	"strings"
	"time"
)

type PlainGenerator struct {

}

func (*PlainGenerator)Generate(content string, level Level) string  {
	fun, file, line := GetCaller()

	str := fmt.Sprintf("%s --【%s】%s %s:%s:%d ",
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ToUpper(level.GetLevelText()),  content,
		file, fun, line, )

	return str
}
