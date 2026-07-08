package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const DEFAULT_TTL = "60" // 1 min TTL

type Command struct {
	op string
	key string
	val string
	ttl string
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

		_, err = reader.ReadString('\n')
		if err != nil {
			return Command{}, err
		}
	}

	return buildCommand(cmd, replyCh)
}

func buildCommand(args []string, replyCh chan Result) (Command, error) {
	switch args[0] {
	case "EXISTS":
		if len(args) != 2 {
			return Command{}, errors.New("improper number of arguments")
		}
		return Command{
			op: "EXISTS",
			key: args[1],
			replyCh: replyCh,
		}, nil

	case "EXPIRE":
		if len(args) > 3 || len(args) < 2{
			return Command{}, errors.New("improper number of arguments")
		}

		exp := DEFAULT_TTL
		if len(args) == 3 {
			exp = args[2]
		}

		return Command{
			op: "EXPIRE",
			key: args[1],
			ttl: exp,
			replyCh: replyCh,
		}, nil

	case "TTL":
		if len(args) != 2 {
			return Command{}, errors.New("improper number of arguments")
		}

		return Command{
			op: "TTL",
			key: args[1],
			replyCh: replyCh,
		}, nil

	case "GET":
		if len(args) != 2 {
			return Command{}, errors.New("improper number of arguments")
		}

		return Command{
			op: "GET",
			key: args[1],
			replyCh: replyCh,
		}, nil

	case "SET":
		if len(args) > 5 || len(args) < 3 {
			return Command{}, errors.New("improper number of arguments")
		}
		
		ttl := "-1"
		if len(args) > 3 {
			if args[3] == "EX" {
				if len(args) == 5 {
					ttl = args[4]
				}

			} else {
				return Command{}, fmt.Errorf("command \"%s\" not recognized", args[3])
			}
		}

		return Command{
			op: "SET",
			key: args[1],
			val: args[2],
			ttl: ttl,
			replyCh: replyCh,
		}, nil

	case "DEL":
		if len(args) != 2 {
			return Command{}, errors.New("improper number of arguments")
		}

		return Command{
			op: "DEL",
			key: args[1],
			replyCh: replyCh,
		}, nil

	default:
		return Command{}, fmt.Errorf("unrecognized command name: %s", args[0])

	}
}

func WriteResp(conn net.Conn, result Result) {
	if result.err != nil {
		msg := "-ERR " + result.err.Error() + "\r\n"
		conn.Write([]byte(msg))
	} else {
		msg := "+" + result.val + "\r\n"
		conn.Write([]byte(msg))
	}
}
