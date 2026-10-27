package server

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/abhikred/go-kv/internal/resp"
	"github.com/abhikred/go-kv/internal/store"
)

const (
	errWrongArgs      = "-ERR wrong number of arguments\r\n"
	errUnknownCommand = "-ERR unknown command\r\n"
	replyOK           = "+OK\r\n"
	replyNil          = "$-1\r\n"
)

const idleTimeout = 5 * time.Minute

func ListenAndServe(ctx context.Context, addr string, s store.Store) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	var wg sync.WaitGroup
	go func() {
		<-ctx.Done()
		ln.Close()
	}()
	var connErr error
	for {
		conn, connAcceptErr := ln.Accept()
		if connAcceptErr != nil {
			connErr = connAcceptErr
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConn(ctx, conn, s)
		}()
	}
	wg.Wait()

	if ctx.Err() != nil {
		return nil
	}
	return connErr
}
func handleConn(ctx context.Context, conn net.Conn, s store.Store) {
	defer conn.Close()
	reader := resp.NewReader(conn)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(idleTimeout)); err != nil {
			return
		}
		parts, err := reader.ReadCommand()
		if err != nil {
			return
		}
		switch strings.ToUpper(parts[0]) {
		case "GET":
			if len(parts) == 2 {
				res, err := s.Get(parts[1])
				if err != nil {
					conn.Write([]byte(replyNil))
				} else {
					fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(res), res)
				}
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		case "SET":
			if len(parts) == 3 {
				s.Set(parts[1], parts[2])
				conn.Write([]byte(replyOK))
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		case "DEL":
			if len(parts) == 2 {
				if err := s.Del(parts[1]); err != nil {
					conn.Write([]byte(":0\r\n"))
				} else {
					conn.Write([]byte(":1\r\n"))
				}
			} else {
				conn.Write([]byte(errWrongArgs))
			}
		default:
			conn.Write([]byte(errUnknownCommand))
		}
	}
}
