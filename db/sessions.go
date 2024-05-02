package db

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var Store = redis.New(redis.Config{
	Host:      "127.0.0.1",
	Port:      6379,
	Username:  "",
	Password:  "",
	Database:  0,
	Reset:     false,
	TLSConfig: nil,
	PoolSize:  10 * runtime.GOMAXPROCS(0),
})

var Sessions = session.New(session.Config{
	Expiration: 1 * time.Hour,
	Storage:    Store,
})
