package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	sqlForCreateCourts = "INSERT INTO courts(id,company_id,points,date,type,courtcode,party1,party2,content) VALUES(?,?,?,?,?,?,?,?,?);"
)

type Court struct {
	ID        int    `orm:"column(id)"`
	CompanyID int    `orm:"column(company_id)"`
	Points    int    `orm:"column(points)"`
	Date      string `orm:"column(date)"`
	Type      string `orm:"column(type)"`
	CourtCode string `orm:"column(court_code)"`
	Party1    string `orm:"column(party1)"`
	Party2    string `orm:"column(party2)"`
	Content   string `orm:"column(content)"`
}

func (c *Court) TableName() string {
	return "courts"
}

func getCourtsKey(cid int) string {
	return fmt.Sprintf("/companies/%d/courts", cid)
}

func getCourtKey(cid, id int) string {
	return fmt.Sprintf("/companies/%d/courts/%d", cid, id)
}

func getCourtsFromRedis(cid int) (*[]Court, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getCourtsKey(cid)))
	if err != nil {
		return nil, err
	}

	var cs []Court

	err = json.Unmarshal(b, &cs)
	return &cs, err
}

func getCourtFromRedis(cid, id int) (*Court, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getCourtKey(cid, id)))
	if err != nil {
		return nil, err
	}

	var c Court

	err = json.Unmarshal(b, &c)
	return &c, err
}

func setCourtsToRedis(cid int, cs *[]Court) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getCourtsKey(cid), b, "EX", EXTime)
	return err
}

func setCourtToRedis(c *Court) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getCourtKey(c.CompanyID, c.ID), b, "EX", EXTime)
	return err
}

func GetCourtsByCompanyID(cid int) (*[]Court, error) {
	cts, err := getCourtsFromRedis(cid)
	if err != nil {
		var cs []Court
		_, err := O.Raw("SELECT * FROM courts WHERE company_id = ?;", cid).QueryRows(&cs)
		if err == nil {
			setCourtsToRedis(cid, &cs)
		}
		return &cs, err
	}
	return cts, err
}

func GetCourtByIDAndCompanyID(id, cid int) (*Court, error) {
	ct, err := getCourtFromRedis(cid, id)
	if err != nil {
		var c Court
		err := O.Raw("SELECT * FROM courts WHERE id = ? AND company_id = ?;", id, cid).QueryRow(&c)
		if err == nil {
			setCourtToRedis(&c)
		}
		return &c, err
	}
	return ct, err
}

func CreateCourt(c *Court) error {
	if c == nil {
		return errors.New("arg c is nil")
	}

	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	var maxID int
	err := o.Raw("SELECT MAX(id) FROM courts WHERE company_id = ?;", c.CompanyID).QueryRow(&maxID)

	_, err = o.Raw(sqlForCreateCourts,
		maxID+1,
		c.CompanyID,
		c.Points,
		c.Date,
		c.Type,
		c.CourtCode,
		c.Party1,
		c.Party2,
		c.Content).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	err = UpdateCompanyPoints(c.CompanyID, -c.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

func DeleteCourtByIDAndCompanyID(id, cid int) error {
	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	c, err := GetCourtByIDAndCompanyID(id, cid)
	if err != nil {
		return err
	}

	err = UpdateCompanyPoints(cid, c.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Raw("DELETE FROM lawsuits WHERE id = ? AND company_id = ?;", id, cid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	return o.Commit()
}
