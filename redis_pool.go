package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

const (
	maxConnections = 5
	connectTimeout = time.Duration(10) * time.Second
	readTimeout = time.Duration(10) * time.Second
	writeTimeout = time.Duration(10) * time.Second
	server = "localhost:6379"
)

var redisPool = redis.NewPool(func() (redis.Conn, error) {
	c, err := redis.DialTimeout("tcp", server, connectTimeout, readTimeout, writeTimeout)
	if err != nil {
		return nil, err
	}

	return c, err
}, maxConnections)