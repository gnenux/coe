package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	sqlForCreateBlackList = "INSERT INTO blacklist(id,company_id,points,reason,date,decision) VALUES(?,?,?,?,?,?);"
)

type BlackList struct {
	ID        int    `orm:"column(id)"`
	CompanyID int    `orm:"column(company_id)"`
	Points    int    `orm:"column(points)"`
	Reason    string `orm:"column(reason)"`
	Date      string `orm:"column(date)"`
	Decision  string `orm:"column(decision)"`
}

func (bl *BlackList) TableName() string {
	return "blacklist"
}

func getBlackListKey(cid, id int) string {
	return fmt.Sprintf("/companies/%d/blacklist/%d", cid, id)
}

func getBlackListsKey(cid int) string {
	return fmt.Sprintf("/companies/%d/blacklist", cid)
}

func getBlackListsFromRedis(cid int) (*[]BlackList, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getBlackListsKey(cid)))
	if err != nil {
		return nil, err
	}

	var bls []BlackList
	err = json.Unmarshal(b, &bls)
	return &bls, err
}

func setBlackListsToRedis(cid int, bls *[]BlackList) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(bls)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getBlackListsKey(cid), b, "EX", EXTime)
	return err
}

func getBlackListFromRedis(cid, id int) (*BlackList, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getBlackListKey(cid, id)))
	if err != nil {
		return nil, err
	}

	var bl BlackList
	err = json.Unmarshal(b, &bl)
	return &bl, err
}

func setBlackListToRedis(bl *BlackList) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(bl)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getBlackListKey(bl.CompanyID, bl.ID), b, "EX", EXTime)
	return err
}

func GetBlackListByCompanyID(cid int) (*[]BlackList, error) {
	bs, err := getBlackListsFromRedis(cid)
	if err != nil {
		var bls []BlackList
		_, err := O.Raw("SELECT * FROM blacklist WHERE company_id = ?;", cid).QueryRows(&bls)
		if err == nil {
			setBlackListsToRedis(cid, &bls)
		}
		return &bls, err
	}
	return bs, err

}

func GetBlackListByIDAndCompanyID(id, cid int) (*BlackList, error) {
	b, err := getBlackListFromRedis(cid, id)
	if err != nil {
		var bl BlackList
		err := O.Raw("SELECT * FROM blacklist WHERE id = ? AND company_id = ?;", id, cid).QueryRow(&bl)
		if err == nil {
			setBlackListToRedis(&bl)
		}
		return &bl, err
	}

	return b, err

}

func CreateBlackList(bl *BlackList) error {
	if bl == nil {
		return errors.New("arg bl is nil")
	}

	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	var maxID int
	err := o.Raw("SELECT MAX(id) FROM blacklist WHERE company_id = ?;", bl.CompanyID).QueryRow(&maxID)

	_, err = o.Raw(sqlForCreateBlackList,
		maxID+1,
		bl.CompanyID,
		bl.Points,
		bl.Reason,
		bl.Date,
		bl.Decision).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	err = UpdateCompanyPoints(bl.CompanyID, -bl.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

func DeleteBlackListByIDAndCompanyID(id, cid int) error {
	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	bl, err := GetBlackListByIDAndCompanyID(id, cid)
	if err != nil {
		return err
	}

	err = UpdateCompanyPoints(cid, bl.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Raw("DELETE FROM blacklist WHERE id = ? AND company_id = ?;", id, cid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	return o.Commit()
}
