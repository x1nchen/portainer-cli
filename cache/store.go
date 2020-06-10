package cache

// Store is the data cache for portainer
type Store interface {
	//
	Save()
}

var _ Store = &bboltStore{}

type bboltStore struct {
	host string // host is for bbolt bucket
}

func NewBoltStore(host string) Store {
	s := &bboltStore{
		host: host,
	}
	return s
}

func (b bboltStore) Save() {
	panic("implement me")
}

