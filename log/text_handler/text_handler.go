package text_handler

const (
	FatalLevel=iota
	ErrorLevel
	DebugLevel
	InfoLevel
)

var Handlers map[string]Handler = make(map[string]Handler, 8)

type Handler interface {
	HandleText(text string, level, minLevel int) error
	LoadHandler(option HandlerOption) error
}
