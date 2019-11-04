package log

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type TextGenerator interface {
	generate(content string, level Level) string
}

func getCaller() (funcName, file string, line int) {

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
