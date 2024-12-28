package actor

import (
	"fmt"
)

type Cb func(pid PID)

var newActorCb Cb
var destroyActorCb Cb

// SetCb  newCb and destroyCb not thread safe!!!
func SetCb(newCb, destroyCb Cb) {
	newActorCb = newCb
	destroyActorCb = destroyCb
}

func RegisterActor(a Actor, cap int) error {
	err := root.spawn(a, cap)
	if err != nil {
		logger.Fatal("actor %v register failed %v", a.PID(), err)
	} else {
		if newActorCb != nil {
			newActorCb(a.PID())
		}
		logger.Info("actor %v register success", a.PID())
	}
	return err
}

func SyncRequest(from, to PID, message *Message) RespMessage {
	ref, ret := root.get(to)
	if !ret {
		return RespMessage{
			Err: fmt.Errorf("from %v request %v failed, target actor nil", from, to),
		}
	}

	message.from = from
	message.to = to
	return ref.Request(message)
}

// AsyncRequest async request, not block, cb will execute in request goroutine
func AsyncRequest(from, to PID, message *Message, cb func(msg RespMessage)) {
	ref, ret := root.get(to)
	if !ret {
		resp := RespMessage{
			Err: fmt.Errorf("from %v request %v failed, target actor nil", from, to),
		}
		cb(resp)
		return
	}
	message.from = from
	message.to = to
	ref.AsyncRequest(message, cb)
}

func Send(from, to PID, message *Message) error {
	ref, ret := root.get(to)
	if !ret {
		return fmt.Errorf("from %v send to %v failed, target actor nil", from, to)
	}
	message.from = from
	message.to = to

	return ref.Send(message)
}

func Stop() {
	root.stop()
}

func StopActor(id PID) {
	ref, ret := root.get(id)
	if !ret {
		return
	}
	ref.Stop()
	root.deleteActor(id)
	logger.Info("actor %v rcv stop", id)
	if destroyActorCb != nil {
		destroyActorCb(id)
	}
}
