package main

import (
	"errors"
	"time"
)

type Store struct {
	vals map[string]string
	ttl map[string]time.Time
}

func NewStore() Store {
	return Store {
		vals: make(map[string]string),
		ttl: map[string]time.Time{},
	}
}

func RunStore(commands chan Command) {
	// TODO: change store from basic hash map to a struct of some kind
	store := NewStore()

	for cmd := range commands {
		switch cmd.op {
		case "GET":
			val, ok := store.vals[cmd.key]
			if ok {
				cmd.replyCh <- Result{
					val: val,
					err: nil,
				}
			} else  {
				cmd.replyCh <- Result{
					val: "nil",
					err: nil,
				}
			} 
		case "EXISTS":
			_, ok := store.vals[cmd.key]
			if ok {
				cmd.replyCh <- Result{
					val: "1",
					err: nil,
				}
			} else {
				cmd.replyCh <- Result{
					val: "0",
					err: nil,
				}
			}
		case "TTL":
		case "SET":
			store.vals[cmd.key] = cmd.val

			cmd.replyCh <- Result{
				val: "OK",
				err: nil,
			}
		case "DEL":
			_, ok := store.vals[cmd.key]
			if ok {
				delete(store.vals, cmd.key)
				delete(store.ttl, cmd.key)

				cmd.replyCh <- Result{
					val: "1",
					err: nil,
				}
			} else {
				cmd.replyCh <- Result{
					val: "0",
					err: nil,
				}
			}
		default:
			cmd.replyCh <- Result{
				val: "nil",
				err: errors.New("Error executing command"),
			}
		}
	}
}
