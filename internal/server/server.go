package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/abhikred/go-kv/internal/store"
)

const errWrongArgs = "ERR wrong number of arguments\n"

func ListenAndServe(addr string, s store.Store) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, s)
	}

}
func handleConn(conn net.Conn, s store.Store) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "GET":
			if len(parts) == 2 {
				res, err := s.Get(parts[1])
				if err != nil {
					conn.Write([]byte(err.Error() + "\n"))
				} else {
					conn.Write([]byte(res + "\n"))
				}
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		case "SET":
			if len(parts) == 3 {
				s.Set(parts[1], parts[2])
				conn.Write([]byte("ok\n"))
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		case "DEL":
			if len(parts) == 2 {
				err := s.Del(parts[1])
				if err != nil {
					conn.Write([]byte(err.Error() + "\n"))
				} else {
					conn.Write([]byte("ok\n"))
				}
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		default:
			res := "Unknown Request"
			conn.Write([]byte(res + "\n"))
		}
	}

}
