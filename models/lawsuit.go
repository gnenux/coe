package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	sqlForCreateLawsuit = "INSERT INTO lawsuits(id,company_id,points,name,type,date) VALUES(?,?,?,?,?,?);"
)

type Lawsuit struct {
	ID        int    `orm:"column(id)"`
	CompanyID int    `orm:"column(company_id)"`
	Points    int    `orm:"column(points)"`
	Name      string `orm:"column(name)"`
	Date      string `orm:"column(date)"`
	Type      string `orm:"column(type)"`
}

func (ls *Lawsuit) TableName() string {
	return "lawsuits"
}

func getLawsuitsKey(cid int) string {
	return fmt.Sprintf("/companies/%d/lawsuits", cid)
}

func getLawsuitKey(cid, id int) string {
	return fmt.Sprintf("/companies/%d/lawsuits/%d", cid, id)
}

func getLawsuitsFromRedis(cid int) (*[]Lawsuit, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getLawsuitsKey(cid)))
	if err != nil {
		return nil, err
	}

	var ls []Lawsuit

	err = json.Unmarshal(b, &ls)
	return &ls, err
}

func getLawsuitFromRedis(cid, id int) (*Lawsuit, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getLawsuitKey(cid, id)))
	if err != nil {
		return nil, err
	}

	var l Lawsuit

	err = json.Unmarshal(b, &l)
	return &l, err
}

func setLawsuitsToRedis(cid int, ls *[]Lawsuit) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(ls)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getLawsuitsKey(cid), b, "EX", EXTime)
	return err
}

func setLawsuitToRedis(l *Lawsuit) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(l)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getLawsuitKey(l.CompanyID, l.ID), b, "EX", EXTime)
	return err
}

func GetLawsuitsByCompanyID(cid int) (*[]Lawsuit, error) {
	ls, err := getLawsuitsFromRedis(cid)
	if err != nil {
		var lss []Lawsuit
		_, err := O.Raw("SELECT * FROM lawsuits WHERE company_id = ?;", cid).QueryRows(&lss)
		if err == nil {
			setLawsuitsToRedis(cid, &lss)
		}
		return &lss, err
	}

	return ls, err
}

func GetLawsuitByIDAndCompanyID(id, cid int) (*Lawsuit, error) {
	l, err := getLawsuitFromRedis(cid, id)
	if err != nil {
		var ls Lawsuit
		err := O.Raw("SELECT * FROM lawsuits WHERE id = ? AND company_id = ?;", id, cid).QueryRow(&ls)
		if err == nil {
			setLawsuitToRedis(&ls)
		}
		return &ls, err
	}

	return l, err
}

func CreateLawsuit(ls *Lawsuit) error {
	if ls == nil {
		return errors.New("arg ls is nil")
	}

	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	var maxID int
	err := o.Raw("SELECT MAX(id) FROM lawsuits WHERE company_id = ?;", ls.CompanyID).QueryRow(&maxID)

	_, err = o.Raw(sqlForCreateLawsuit,
		maxID+1,
		ls.CompanyID,
		ls.Points,
		ls.Name,
		ls.Type,
		ls.Date).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	err = UpdateCompanyPoints(ls.CompanyID, -ls.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

func DeleteLawsuitByIDAndCompanyID(id, cid int) error {
	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	ls, err := GetLawsuitByIDAndCompanyID(id, cid)
	if err != nil {
		return err
	}

	err = UpdateCompanyPoints(cid, ls.Points, o)
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
