package log

type Handler interface {
	HandleText(text string)
	Load() error
}

