package redis

func GetRedisConfig(config map[string]string) *RedisConfig {
	return &RedisConfig{
		Addr:     config["addr"],
		Username: config["username"],
		Password: config["password"],
	}
}

func GetRedisConnectorConfig(config map[string]string) *ConnectorConfig {
	var mode ConnectorMode
	modeStr := config["mode"]
	switch modeStr {
	default:
		mode = Standalone
	}
	return &ConnectorConfig{
		Mode:        mode,
		RedisConfig: GetRedisConfig(config),
	}
}
