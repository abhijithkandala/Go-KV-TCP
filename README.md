# Go-KV-TCP
Redis-compatible in-memory key-value store built from scratch in Go TCP, RESP2, snapshots, TTL.
This is a learning project, a prerequisite for my main project 'KV store with Raft concensus'.


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