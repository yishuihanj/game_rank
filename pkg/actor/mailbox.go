package actor

const minQueueSize = 32

type MailBox interface {
	Enqueue(m interface{}) error
	OutCh() <-chan interface{}
}

type defaultMailBox struct {
	outCh chan interface{}
	cap   int
}

func newDefaultMailBox(cap int) *defaultMailBox {
	if cap < minQueueSize {
		cap = minQueueSize
	}
	m := &defaultMailBox{
		outCh: make(chan interface{}, cap),
		cap:   cap,
	}

	return m
}

func (d *defaultMailBox) Enqueue(m interface{}) error {
	if len(d.outCh) >= d.cap { // buf overflow
		return ErrMailOverflow
	}
	d.outCh <- m
	return nil
}

func (d *defaultMailBox) OutCh() <-chan interface{} {
	return d.outCh
}
