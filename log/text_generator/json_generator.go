package text_generator

import (
	"encoding/json"
	"time"
)

type JsonGenerator struct {
}

type JsonInfo struct {
	CallerFunc string    `json:"func"`
	CallerFile string    `json:"file"`
	CallerLine int       `json:"line"`
	Level      string    `json:"level"`
	Time       time.Time `json:"time"`
	Tag        string    `json:"tag"`
	Content    string    `json:"content"`
}

func init()  {
	TextGenerators["json"] = new(JsonGenerator)
}

func (*JsonGenerator) Generate(tag, content, level string) string {
	fun, file, line := GetCaller()
	j := JsonInfo{fun,file,line,level,time.Now(),tag,content}

	b, err := json.Marshal(j)
	if err != nil {
		return ""
	}

	return string(b)
}
