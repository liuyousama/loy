package text_generator

import "runtime"

var TextGenerators map[string]TextGenerator = make(map[string]TextGenerator, 8)

type TextGenerator interface {
	Generate(tag, content, level string) string
}

func GetCaller() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "未知文件(unknown file)"
		line = 0
		funcName = "未知调用方法(unknown caller func)"
	}

	f := runtime.FuncForPC(pc)
	funcName = f.Name()
	return
}
