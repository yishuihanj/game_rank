package actor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	criticalTIme = time.Millisecond * 100
)

type Handler func(ctx context.Context, data *Message) error

type StatisticFunc func(uid uint64, id uint32, ms int64)

type Dispatcher struct {
	handlers       map[uint32]Handler
	maxExecuteTime time.Duration
	st             StatisticFunc
	sync.RWMutex
}

func NewDispatcher(maxMillSec int, st StatisticFunc) *Dispatcher {
	return &Dispatcher{
		handlers:       make(map[uint32]Handler),
		st:             st,
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
	// unregister id
	h, ok := d.handlers[id]
	if !ok {
		err := fmt.Errorf("from %v to %v uid: %d rcv not registered id: %d", data.from, data.to, uid, id)
		logger.Fatal(err.Error())
		return err
	}

	// handle
	bt := time.Now()
	err := h(ctx, data)
	et := time.Since(bt)
	if err != nil {
		logger.Error("dispatch failed: from %v to %v uid %v id %v, err: %v", data.from, data.to, uid, id, err)
	}

	if et > d.maxExecuteTime {
		if et > criticalTIme {

			logger.Fatal("handle slow, from %v to %v uid %v id %v cost %v ms", data.from, data.to, uid, id,
				et.Milliseconds())
		} else {
			logger.Error("handle slow, from %v to %v uid %v id %v cost %v ms", data.from, data.to, uid, id,
				et.Milliseconds())
		}
	}

	if d.st != nil {
		d.st(uid, id, et.Milliseconds())
	}

	return err
}
