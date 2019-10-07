package text_handler

import (
	"fmt"
	"os"
)

func init()  {
	Handlers["console"] = new(ConsoleHandler)
}

type ConsoleHandler struct {

}

func (h *ConsoleHandler)LoadHandler(option HandlerOption) error {
	return nil
}

func (*ConsoleHandler)HandleText(text string, level, minLevel int) {
	if level > minLevel {
		return
	}

	retryExecutor(func() error {
		_, err := fmt.Fprintln(os.Stdout, text)
		return err
	})
}
