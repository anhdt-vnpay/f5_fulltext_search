package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl/db_connector"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/log"
	"github.com/anhdt-vnpay/f5_fulltext_search/runtime"
	"github.com/anhdt-vnpay/f5_fulltext_search/search_processor"
	"github.com/golang/glog"

	gorm_impl "github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl"
	redis_pubsub "github.com/anhdt-vnpay/f5_fulltext_search/message_processor"
	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
	service "github.com/anhdt-vnpay/f5_fulltext_search/service/gorm"
)

var (
	REQUEST_TIMEOUT = 5 * time.Second
	apiLogger       = log.NewLogger("api.gateway")
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Range")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Next
		next.ServeHTTP(w, r)
		return
	})
}

func httpHandlers(listener net.Listener) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	// -------------------------
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
		apiLogger.Errorf("Create database connector error: %s", err.Error())
		return err
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

	ch := messageProcessor.GetMsgChannel()
	searchProcessor := search_processor.NewSearchProcessor(ch, redisChannel, redisConnector)

	var opts []runtime.DbOption
	opt1 := runtime.WithStorage(dbStorage)
	opt2 := runtime.WithMsgProcessor(messageProcessor)
	opt3 := runtime.WithSearchProcessor(searchProcessor)
	opts = append(opts, opt1, opt2, opt3)
	dbfs := runtime.NewDbFullTextSearch(opts...)

	userManager := service.NewUserManager(dbfs)
	// -------------------------

	router := http.NewServeMux()
	router.HandleFunc("/es/get", userManager.Get)
	router.HandleFunc("/es/insert", userManager.Insert)
	router.HandleFunc("/es/update", userManager.Update)
	router.HandleFunc("/es/delete", userManager.Delete)
	router.HandleFunc("/es/search", userManager.SearchLite)

	cors_handler := CORS(router)
	server := &http.Server{
		ReadTimeout: REQUEST_TIMEOUT,
		Handler:     cors_handler,
	}

	return server.Serve(listener)
}

func GatewayServer(port int) error {
	// Load config

	config.Init()

	defer glog.Flush()

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		msg := fmt.Sprintf("net.Listener error: %s", err)
		apiLogger.Errorf("Serve httpHandlers error: %s", msg)
		return err
	}

	defer listener.Close()

	apiLogger.Infof("Server start listening in %s", listener.Addr().String())

	if err := httpHandlers(listener); err != nil {
		glog.Fatal(err)
	}

	return nil
}
