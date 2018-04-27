package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	sqlForCreatePenalty = "INSERT INTO penalties(id,company_id,points,punish_number,content,decision,Date) VALUES(?,?,?,?,?,?,?);"
)

type Penalty struct {
	ID           int    `orm:"column(id)"`
	CompanyID    int    `orm:"column(company_id)"`
	Points       int    `orm:"column(points)"`
	PunishNumber string `orm:"column(punish_number)"`
	Content      string `orm:"column(content)"`
	Decision     string `orm:"column(decision)"`
	Date         string `orm:"column(date)"`
}

func (p *Penalty) TableName() string {
	return "penalties"
}

func getPenaltiesKey(cid int) string {
	return fmt.Sprintf("/companies/%d/enalties", cid)
}

func getPenaltyKey(cid, id int) string {
	return fmt.Sprintf("/companies/%d/penalties/%d", cid, id)
}

func getPenaltiesFromRedis(cid int) (*[]Penalty, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getPenaltiesKey(cid)))
	if err != nil {
		return nil, err
	}

	var ps []Penalty

	err = json.Unmarshal(b, &ps)
	return &ps, err
}

func getPenaltyFromRedis(cid, id int) (*Penalty, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getPenaltyKey(cid, id)))
	if err != nil {
		return nil, err
	}

	var p Penalty

	err = json.Unmarshal(b, &p)
	return &p, err
}

func setPenaltiesToRedis(cid int, ps *[]Penalty) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getPenaltiesKey(cid), b, "EX", EXTime)
	return err
}

func setPenaltyToRedis(p *Penalty) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getPenaltyKey(p.CompanyID, p.ID), b, "EX", EXTime)
	return err
}

func GetPenaltiesByCompanyID(cid int) (*[]Penalty, error) {
	pes, err := getPenaltiesFromRedis(cid)
	if err != nil {
		var ps []Penalty
		_, err := O.Raw("SELECT * FROM penalties WHERE company_id = ?;", cid).QueryRows(&ps)
		if err == nil {
			setPenaltiesToRedis(cid, &ps)
		}
		return &ps, err
	}

	return pes, err
}

func GetPenaltyByIDAndCompanyID(id, cid int) (*Penalty, error) {
	pe, err := getPenaltyFromRedis(cid, id)
	if err != nil {
		var p Penalty
		err := O.Raw("SELECT * FROM penalties WHERE id = ? AND company_id = ?;", id, cid).QueryRow(&p)
		if err == nil {
			setPenaltyToRedis(&p)
		}
		return &p, err
	}

	return pe, err
}

func CreatePenalty(p *Penalty) error {
	if p == nil {
		return errors.New("arg p is nil")
	}

	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	var maxID int
	err := o.Raw("SELECT MAX(id) FROM penalties WHERE company_id = ?;", p.CompanyID).QueryRow(&maxID)

	_, err = o.Raw(sqlForCreatePenalty,
		maxID+1,
		p.CompanyID,
		p.Points,
		p.PunishNumber,
		p.Content,
		p.Decision,
		p.Date).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	err = UpdateCompanyPoints(p.CompanyID, -p.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

func DeletePenaltyByIDAndCompanyID(id, cid int) error {
	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}
	p, err := GetPenaltyByIDAndCompanyID(id, cid)
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Raw("DELETE FROM penalties WHERE id = ? AND company_id = ?;", id, cid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	err = UpdateCompanyPoints(cid, p.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}
