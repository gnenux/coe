package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
)

var O orm.Ormer

var RedisPool *redis.Pool

const (
	EXTime = "3600"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/coe?charset=utf8")
	O = orm.NewOrm()

	RedisPool = newPool(":6379")
}

func NewOrm() orm.Ormer {
	return orm.NewOrm()
}

func UpdateCompanyPoints(id, points int, o orm.Ormer) error {
	originalPoints := 100
	err := o.Raw("SELECT points FROM companies WHERE id = ?;", id).QueryRow(&originalPoints)
	if err != nil {
		return err
	}

	originalPoints = originalPoints + points

	_, err = o.Raw("UPDATE companies SET points = ? WHERE id = ?;", originalPoints, id).Exec()
	return err
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     12,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}
