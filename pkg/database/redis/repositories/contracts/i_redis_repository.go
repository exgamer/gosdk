package contracts

type IRedisRepository interface {
	FlushAll() error
}
