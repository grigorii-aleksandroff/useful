package locker

type StoreType int

const (
	Memory StoreType = iota
	Redis
	Mock // position can be changed
)
