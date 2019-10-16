package tcp

const _jogWorkPoolSize = 16

type worker struct {
	jobChan chan *request
	router Router
}

func (w *worker)work()  {
	for {
		select {
		case request := <- w.jobChan:
			h, ok := w.router.handlers[request.msg.pathId]
			if !ok {
				continue
			}
			h=h
		}
	}
}