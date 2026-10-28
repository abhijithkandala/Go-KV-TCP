# Go-KV-TCP
Redis-compatible in-memory key-value store built from scratch in Go TCP, RESP2, snapshots, TTL.
This is a learning project, a prerequisite for my main project 'KV store with Raft concensus'.


## Running it

Requires Go 1.27+. The server listens on `127.0.0.1:6380` (hardcoded in
`main.go`), one port above real Redis so both can run side by side.

```bash
git clone https://github.com/abhikred/go-kv
cd go-kv
go run .
```

Press Ctrl+C to stop the listener closes, in-flight connections drain,
then the store's goroutine shuts down.

### Talking to it

It speaks RESP2, so `redis-cli` works as-is:

```bash
redis-cli -p 6380
127.0.0.1:6380> SET foo bar
OK
127.0.0.1:6380> GET foo
"bar"
127.0.0.1:6380> DEL foo
(integer) 1
```

Supported commands: `GET key`, `SET key value`, `DEL key`.

Plain `nc 127.0.0.1 6380` will not work for typing commands by hand  the
parser only accepts RESP arrays, so a bare `SET foo bar` is rejected and the
connection closes. To use `nc`, send the wire format:

```bash
printf '*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n' | nc 127.0.0.1 6380
```

### Tests

```bash
go test -race ./...
```


## Purpose of each component

1. **store** — the key-value logic. Two implementations behind one `Store`
   interface: `ActorStore` owns the map on its own goroutine and callers send
   requests over a channel, so writes are serialized without a lock;
   `MutexStore` guards the same map with a `sync.RWMutex`.
2. **server** — the TCP listener and connection handler. Accepts connections,
   serves each on its own goroutine, dispatches parsed commands to the store
   and writes the replies. Also handles graceful shutdown: on SIGINT it closes
   the listener and drains in-flight connections, with an idle read deadline
   so a silent client can't hold it up.
3. **resp** — decodes the RESP wire format (REdis Serialization Protocol) into
   a command and its arguments. Protocol only: it returns a `[]string` and
   knows nothing about GET, SET or the store.