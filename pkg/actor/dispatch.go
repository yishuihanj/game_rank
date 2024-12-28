package actor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Handler func(ctx context.Context, data *Message) error
type Dispatcher struct {
	handlers       map[uint32]Handler
	maxExecuteTime time.Duration
	sync.RWMutex
}

func NewDispatcher(maxMillSec int) *Dispatcher {
	return &Dispatcher{
		handlers:       make(map[uint32]Handler),
		maxExecuteTime: time.Duration(maxMillSec) * time.Millisecond,
	}
}

func (d *Dispatcher) Register(id uint32, cmd Handler) {
	if _, ok := d.handlers[id]; ok {
		panic(fmt.Sprintf("register id %d duplicated!", id))
	}
	d.handlers[id] = cmd
}
func (d *Dispatcher) Dispatch(ctx context.Context, data *Message) error {
	id := data.Id
	uid := data.Uid
	h, ok := d.handlers[id]
	if !ok {
		err := fmt.Errorf("from %v to %v uid: %d rcv not registered id: %d", data.from, data.to, uid, id)
		logger.Fatal(err.Error())
		return err
	}
	err := h(ctx, data)
	return err
}
