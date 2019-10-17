package tcp

import "sync"

type connectManager struct {
	connectMap map[interface{}]*connect
	mu sync.RWMutex
}

var cm = &connectManager{
	connectMap:make(map[interface{}]*connect, 0),
}

func saveConnect(key interface{}, conn *connect) {
	cm.saveConnect(key, conn)
}
func (c *connectManager)saveConnect(key interface{}, conn *connect) {
	c.mu.Lock()
	c.connectMap[key] = conn
	c.mu.Unlock()
}

func getConnect(key interface{}) *connect {
	return cm.getConnect(key)
}
func (c *connectManager)getConnect(key interface{}) *connect {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, ok := c.connectMap[key]
	if !ok {
		return nil
	}
	return conn
}

func delConnect(key interface{})  {
	cm.delConnect(key)
}
func (c *connectManager)delConnect(key interface{})  {
	c.mu.RLock()

	_, ok := c.connectMap[key]
	if ok {
		delete(c.connectMap, key)
	}

	c.mu.RUnlock()
}


func getAllConnect() []*connect {
	return cm.getAllConnect()
}
func (c *connectManager)getAllConnect() []*connect {
	c.mu.RLock()

	var connects []*connect
	for _, val := range c.connectMap {
		connects = append(connects, val)
	}

	c.mu.RUnlock()
	return connects
}