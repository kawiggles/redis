package main

import (
	"errors"
	"net"
	"bufio"
	"strings"
	"strconv"
	"io"
)

type Command struct {
	op string
	key string
	val string
	replyCh chan Result
}

type Result struct {
	val string
	err error
}

func ParseResp(reader *bufio.Reader, replyCh chan Result) (Command, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Command{}, err
	}

	if line[0] != '*' {
		return Command{}, errors.New("bad request prefix")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return Command{}, err
	}

	cmd := make([]string, count)
	for i := range count {
		line, err = reader.ReadString('\n')
		if err != nil {
			return Command{}, err
		}
		
		length, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return Command{}, err
		}

		buf := make([]byte, length)
		io.ReadFull(reader, buf)
		cmd[i] = string(buf)

		reader.ReadString('\n')
	}

	return buildCommand(cmd, replyCh)
}

func buildCommand(args []string, replyCh chan Result) (Command, error) {
	switch args[0] {
	case "EXISTS":
		if len(args) < 2 {
			return Command{}, errors.New("not enough arguments for EXISTS request")
		}
		return Command{
			op: "EXISTS",
			key: args[1],
			val: "nil",
			replyCh: replyCh,
		}, nil
	case "TTL":
		if len(args) < 2 {
			return Command{}, errors.New("not enough arguments for TTL request")
		}
		return Command{
			op: "TTL",
			key: args[1],
			val: "nil",
			replyCh: replyCh,
		}, nil
	case "GET":
		if len(args) < 2 {
			return Command{}, errors.New("not enough arguments for GET request")
		}
		return Command{
			op: "GET",
			key: args[1],
			val: "nil",
			replyCh: replyCh,
		}, nil
	case "SET":
		if len(args) < 3 {
			return Command{}, errors.New("not enough arguments for SET request")
		}
		return Command{
			op: "SET",
			key: args[1],
			val: args[2],
			replyCh: replyCh,
		}, nil
	case "DEL":
		if len(args) < 3 {
			return Command{}, errors.New("not enough arguments for DEL request")
		}
		return Command{
			op: "DEL",
			key: args[1],
			val: args[2],
			replyCh: replyCh,
		}, nil
	default:
		msg := "unrecognized command name: " + args[0]
		return Command{}, errors.New(msg)
	}
}

func WriteResp(conn net.Conn, result Result) {
	if result.err != nil {
		msg := "-ERR " + result.err.Error() + "/r/n"
		conn.Write([]byte(msg))
	} else {
		msg := "+" + result.val + "/r/n"
		conn.Write([]byte(msg))
	}
}
