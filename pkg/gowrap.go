package util

import (
	"log"
	"runtime/debug"
)

var recoverHandler = RecoverHandler

func SetRecoverHandler(h func(interface{})) {
	recoverHandler = h
}

func Go(f func()) {
	go func() {
		if recoverHandler != nil {
			defer func() {
				if p := recover(); p != nil {
					recoverHandler(p)
				}
			}()
		}
		f()
	}()
}

func RecoverHandler(e interface{}) {
	trace := string(debug.Stack())
	debug.PrintStack()
	log.Println("recover:", e, "\n", trace)
}
