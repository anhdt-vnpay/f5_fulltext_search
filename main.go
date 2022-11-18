package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	gorm_impl "github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl"
	db_connector "github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl/db_connector"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	redis_pubsub "github.com/anhdt-vnpay/f5_fulltext_search/message_processor"
	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
	"github.com/anhdt-vnpay/f5_fulltext_search/runtime"
	"github.com/anhdt-vnpay/f5_fulltext_search/search_processor"
)

var (
	wg sync.WaitGroup
)

func main() {
	fmt.Println("====================== Init DB ==========================")
	config.Init()
	fmt.Println("init mysql database ...")
	config := config.GetConfig()

	// Declare variable configuration mysql
	host := config.GetString("mysql.host")
	port := config.GetString("mysql.port")
	dbName := config.GetString("mysql.database")
	user := config.GetString("mysql.username")
	pass := config.GetString("mysql.password")

	configObj := db_connector.ConnectorConfig{
		Mode:     db_connector.Mysql,
		Host:     host,
		Port:     port,
		Database: dbName,
		Username: user,
		Password: pass,
	}

	conn, err := db_connector.NewDatabaseConnector(&configObj)
	if err != nil {
		log.Fatal("DB connector error: ", err.Error())
	}

	db := conn.GetDb()

	// Init redis
	rAddr := config.GetString("redis-config.addr")
	rUsername := config.GetString("redis-config.username")
	rPassword := config.GetString("redis-config.password")

	redisConfig := redis_connector.RedisConfig{
		Addr:     rAddr,
		Username: rUsername,
		Password: rPassword,
	}

	redisConnectorConfig := redis_connector.ConnectorConfig{
		Mode:        redis_connector.Standalone,
		RedisConfig: &redisConfig,
	}

	redisChannel := config.GetString("redis-config.channelname")
	redisConnector := redis_connector.NewRedisConnector(&redisConnectorConfig)

	fmt.Println("====================== Init objects ==========================")

	dbStorage := gorm_impl.NewGormDbStorage(db)
	messageProcessor := redis_pubsub.NewMessageProcessor(redisChannel, redisConnector)

	ch := make(chan string)
	searchProcessor := search_processor.NewSearchProcessor(ch, redisChannel, redisConnector)
	go func() {
		searchProcessor.Subscribe()
	}()
	go func() {
		searchProcessor.PassMessage()
	}()

	var opts []runtime.DbOption
	opt1 := runtime.WithStorage(dbStorage)
	opt2 := runtime.WithMsgProcessor(messageProcessor)
	opt3 := runtime.WithSearchProcessor(searchProcessor)
	opts = append(opts, opt1, opt2, opt3)
	dbf := runtime.NewDbFullTextSearch(opts...)

	fmt.Println("====================== DEMO ==========================")

	// Test get
	fmt.Println("Update after 3 seconds >>>>>>>>")
	rs := m.User{
		ID:   1,
		Name: "A",
		Type: "Person B",
	}
	time.Sleep(4 * time.Second)
	err = dbf.Update("users", &rs)
	if err != nil {
		fmt.Println("Get error: ", err.Error())
	}

	rs2 := m.User{
		ID:   1,
		Name: "B",
		Type: "Person B",
	}
	time.Sleep(2 * time.Second)
	err = dbf.Update("users", &rs2)
	if err != nil {
		fmt.Println("Get error: ", err.Error())
	}
	fmt.Println("RS: ")
	fmt.Println(rs)
	wg.Add(1)
	wg.Wait()
}
