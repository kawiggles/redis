package redis

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
	for i := 0; i < count; i++ {
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
		return Command{}, errors.New("unrecognized command name")
	}
}

func WriteResp(conn net.Conn, result Result) {
}
