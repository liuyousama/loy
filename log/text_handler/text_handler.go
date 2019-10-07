package text_handler

const (
	FatalLevel=iota
	ErrorLevel
	DebugLevel
	InfoLevel
)

var Handlers map[string]Handler = make(map[string]Handler, 8)

type Handler interface {
	HandleText(text string, level, minLevel int)
	LoadHandler(option HandlerOption) error
}

func retryExecutor(f func() error) {
	var err error
	for i := 0; i < maxRetryTimes; i++ {
		err = f()
		if err != nil {
			continue
		} else {
			break
		}
	}
}