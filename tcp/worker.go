package tcp

const _jobWorkPoolSize = 16

type worker struct {
	jobChan chan *request
	router *Router
}

func newWorker(router *Router) *worker {
	w := &worker{
		jobChan:make(chan *request, 128),
		router:router,
	}

	for i:=0; i<_jobWorkPoolSize; i++ {
		go w.work()
	}

	return w
}

func (w *worker)getAddRequestFunc() func(req *request) {
	return func(req *request) {
		select {
		case w.jobChan <- req:
			return
		}
	}
}

func (w *worker)work()  {
	for {
		select {
		case request := <- w.jobChan:
			h, ok := w.router.handlers[request.msg.pathId]
			if !ok {
				continue
			}

			ctx := newContext(request)
			h(ctx)
		}
	}
}