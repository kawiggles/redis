package redis

import ()

func RunStore(commands chan Command) {
	// TODO: change store from basic hash map to a struct of some kind
	store := make(map[string]string)

	for cmd := range commands {
		switch cmd.op {
		// TODO: Execute basic commands
		// and then write to reply channel in that command
		}
	}
}
