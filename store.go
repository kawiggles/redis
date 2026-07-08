package main

import (
	"errors"
	"time"
)

type Store struct {
	vals map[string]string
	exp map[string]time.Time
}

func NewStore() Store {
	return Store {
		vals: make(map[string]string),
		exp: make(map[string]time.Time),
	}
}

func RunStore(commands chan Command) {
	// TODO: change store from basic hash map to a struct of some kind
	store := NewStore()

	for cmd := range commands {
		switch cmd.op {
		case "GET": store.get(cmd)
		case "EXISTS": store.exists(cmd)
		case "EXPIRE": store.expire(cmd)
		case "TTL": store.ttl(cmd)
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
				delete(store.exp, cmd.key)

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

func (s Store) get(cmd Command) {
	val, ok := s.vals[cmd.key]
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
}

func (s Store) exists(cmd Command) {
	_, ok := s.vals[cmd.key]

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
}

func (s Store) expire(cmd Command) {
	_, ok := s.exp[cmd.key]
	if !ok {
		cmd.replyCh <- Result{
			val: "0",
			err: nil,
		}
	}

	dur, err := time.ParseDuration(cmd.ttl)
	if err != nil {
		cmd.replyCh <- Result{
			val: "0",
			err: err,
		}
	}

	s.exp[cmd.key] = time.Now().Add(dur)
	cmd.replyCh <- Result{
		val: "1",
		err: nil,
	}
}

func (s Store) ttl(cmd Command) {
	_, ok := s.vals[cmd.key]

	if !ok {
		cmd.replyCh <- Result{
			val: "-2",
			err: nil,
		}

	}

	ttl, ok := s.exp[cmd.key]
	if ok {
		if time.Now().Compare(ttl) >= 0 {
			delete(s.vals, cmd.key)
			delete(s.exp, cmd.key)

			cmd.replyCh <- Result{
				val: "-2",
				err: nil,
			}

		} else {
			time := time.Until(ttl).Round(time.Second).String()

			cmd.replyCh <- Result{
				val: time,
				err: nil,
			}

		}
	}
}
