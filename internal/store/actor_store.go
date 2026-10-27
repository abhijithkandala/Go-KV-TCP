package store

type reply struct {
	value string
	err   error
}
type request struct {
	op    string
	key   string
	value string
	reply chan reply
}
type ActorStore struct {
	request chan request
	done    chan struct{}
}

func NewActorStore() *ActorStore {
	ch := make(chan request)
	done := make(chan struct{})
	ac := &ActorStore{request: ch, done: done}
	go run(ch, done)
	return ac
}

func run(requests chan request, done chan struct{}) {
	mp := make(map[string]string)
	for {
		select {
		case req := <-requests:
			switch req.op {
			case "get":
				val, ok := mp[req.key]
				if !ok {
					req.reply <- reply{err: ErrEntryNotFound}
				} else {
					req.reply <- reply{value: val}
				}
			case "set":
				mp[req.key] = req.value
				req.reply <- reply{}
			case "del":
				_, ok := mp[req.key]
				if !ok {
					req.reply <- reply{err: ErrEntryNotFound}
				} else {
					delete(mp, req.key)
					req.reply <- reply{}
				}
			default:
				req.reply <- reply{err: ErrUnknownOperation}
			}
		case <-done:
			return
		}
	}
}
func (as *ActorStore) Get(key string) (string, error) {
	rc := make(chan reply, 1)
	req := request{op: "get", key: key, reply: rc}
	as.request <- req
	rep := <-rc
	return rep.value, rep.err
}

func (as *ActorStore) Set(key, value string) {
	rc := make(chan reply, 1)
	req := request{op: "set", key: key, value: value, reply: rc}
	as.request <- req
	<-rc
}

func (as *ActorStore) Del(key string) error {
	rc := make(chan reply, 1)
	req := request{op: "del", key: key, reply: rc}
	as.request <- req
	rep := <-rc
	return rep.err
}
func (as *ActorStore) Close() {
	close(as.done)
}

var _ Store = (*ActorStore)(nil)
