package tcp

type Router struct {
	handlers map[uint32]Handler
}

type Handler func(ctx *Context)

func NewRouter() *Router {
	return &Router{
		handlers:make(map[uint32]Handler, 0),
	}
}

func (r *Router)Add(pathId uint32, handler Handler) {
	r.handlers[pathId] = handler
}