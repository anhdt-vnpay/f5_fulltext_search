package redis

type ConnectorMode int64

const (
	Standalone ConnectorMode = 0
)

type ConnectorConfig struct {
	Mode        ConnectorMode
	RedisConfig *RedisConfig
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
}
