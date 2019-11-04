package log

type Handler interface {
	handleText(text string)
	load() error
}

