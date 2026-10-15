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
}

func NewActorStore() *ActorStore {
	ch := make(chan request)
	ac := &ActorStore{ch}
	go run(ch)
	return ac
}

func run(requests chan request) {
	mp := make(map[string]string)
	for req := range requests {
	}

}
