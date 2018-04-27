package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	sqlForCreateBadManagement = "INSERT INTO badmanagements(id,company_id,points,reason,date,decision) VALUES(?,?,?,?,?,?);"
)

type BadManagement struct {
	ID        int    `orm:"column(id)"`
	CompanyID int    `orm:"column(company_id)"`
	Points    int    `orm:"column(points)"`
	Reason    string `orm:"column(reason)"`
	Date      string `orm:"column(date)"`
	Decision  string `orm:"column(decision)"`
}

func (bm *BadManagement) TableName() string {
	return "badmanagements"
}

func getBadManagementsKey(cid int) string {
	return fmt.Sprintf("/companies/%d/badmanagements", cid)
}

func getBadManagementKey(cid, id int) string {
	return fmt.Sprintf("/companies/%d/badmanagements/%d", cid, id)
}

func getBadManagementsFromRedis(cid int) (*[]BadManagement, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getBadManagementsKey(cid)))
	if err != nil {
		return nil, err
	}

	var bms []BadManagement
	err = json.Unmarshal(b, &bms)
	return &bms, err

}

func setBadManagementsToRedis(cid int, bms *[]BadManagement) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(bms)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getBadManagementsKey(cid), b, "EX", EXTime)
	return err
}

func getBadManagementFromRedis(cid, id int) (*BadManagement, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getBadManagementKey(cid, id)))
	if err != nil {
		return nil, err
	}

	var bm BadManagement
	err = json.Unmarshal(b, &bm)
	return &bm, err

}

func setBadManagementToRedis(bm *BadManagement) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(bm)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getBadManagementKey(bm.CompanyID, bm.ID), b, "EX", EXTime)
	return err
}

func GetBadManagementsByCompanyID(cid int) (*[]BadManagement, error) {
	bdms, err := getBadManagementsFromRedis(cid)
	if err != nil {
		var bms []BadManagement
		_, err := O.Raw("SELECT * FROM badmanagements WHERE company_id = ?;", cid).QueryRows(&bms)
		if err == nil {
			setBadManagementsToRedis(cid, &bms)
		}
		return &bms, err
	}

	return bdms, err

}

func GetBadManagementByIDAndCompanyID(id, cid int) (*BadManagement, error) {
	bdm, err := getBadManagementFromRedis(cid, id)
	if err != nil {
		var bm BadManagement
		err := O.Raw("SELECT * FROM badmanagements WHERE id = ? AND company_id = ?;", id, cid).QueryRow(&bm)
		if err == nil {
			setBadManagementToRedis(&bm)
		}
		return &bm, err
	}
	return bdm, err
}

func CreateBadManagement(bm *BadManagement) error {
	if bm == nil {
		return errors.New("arg bm is nil")
	}

	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}
	var maxID int
	err := o.Raw("SELECT MAX(id) FROM badmanagements WHERE company_id = ?;", bm.CompanyID).QueryRow(&maxID)

	_, err = o.Raw(sqlForCreateBadManagement,
		maxID+1,
		bm.CompanyID,
		bm.Points,
		bm.Reason,
		bm.Date,
		bm.Decision).Exec()

	if err != nil {
		o.Rollback()
		return err
	}

	err = UpdateCompanyPoints(bm.CompanyID, -bm.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}

func DeleteBadManagementByIDAndCompanyID(id, cid int) error {
	o := NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}

	bm, err := GetBadManagementByIDAndCompanyID(id, cid)
	if err != nil {
		return err
	}

	_, err = o.Raw("DELETE FROM badmanagements WHERE id = ? AND company_id = ?;", id, cid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	err = UpdateCompanyPoints(cid, bm.Points, o)
	if err != nil {
		o.Rollback()
		return err
	}

	return o.Commit()
}
