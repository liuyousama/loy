package log

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var TextGenerators map[string]TextGenerator = make(map[string]TextGenerator, 8)

type TextGenerator interface {
	Generate(content string, level Level) string
}

func GetCaller() (funcName, file string, line int) {

	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "未知文件(unknown file)"
		line = 0
		funcName = "未知调用方法(unknown caller func)"
	}


	base, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file = strings.TrimPrefix(file, base)

	f := runtime.FuncForPC(pc)
	funcName = f.Name()
	return
}
