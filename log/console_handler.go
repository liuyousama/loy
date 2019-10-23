package log

import (
	"fmt"
	"os"
)

//type check
var _ Handler = &ConsoleHandler{}

type ConsoleHandler struct {
	colorOn bool
}

func (h *ConsoleHandler)Load() error {
	return nil
}

func (*ConsoleHandler)HandleText(text string) {
	retryExecutor(func() error {
		_, err := fmt.Fprintln(os.Stdout, text)
		return err
	})
}