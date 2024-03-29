package log

import (
	"encoding/json"
	"time"
)

var _ TextGenerator = &JsonGenerator{}

type JsonGenerator struct {
}

type JsonInfo struct {
	CallerFunc string    `json:"func"`
	CallerFile string    `json:"file"`
	CallerLine int       `json:"line"`
	Level      string    `json:"level"`
	Time       time.Time `json:"time"`
	Content    string    `json:"content"`
}

func NewJsonGenerator() *JsonGenerator {
	return new(JsonGenerator)
}

func (*JsonGenerator) generate(content string, level Level) string {
	fun, file, line := getCaller()
	j := JsonInfo{fun,file,line,level.GetLevelText(),time.Now(), content}

	b, err := json.Marshal(j)
	if err != nil {
		return ""
	}

	return string(b)
}
