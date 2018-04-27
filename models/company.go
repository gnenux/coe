package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gomodule/redigo/redis"
)

type Company struct {
	ID                           int    `orm:"column(id)" form:"-" json:"id"`
	Name                         string `orm:"column(name)" form:"name" json:"name"`
	State                        string `orm:"column(state)" form:"state" json:"state"`
	Points                       int    `orm:"column(points)" form:"points" json:"points"`
	LegalPerson                  string `orm:"column(legal_person)" form:"legal_person" json:"legal_person"`
	RegisteredCapital            string `orm:"column(registered_capital)" form:"registered_capital" json:"registered_capital"`
	RegistrationTime             string `orm:"column(registration_time)" form:"registration_time" json:"registration_time"`
	ApprovalTime                 string `orm:"column(approval_time)" form:"approval_time" json:"approval_time"`
	BusinessRegistrationNumber   string `orm:"column(business_registration_number)" form:"business_registration_number" json:"business_registration_number"`
	OrganizationCode             string `orm:"column(organization_code)" form:"organization_code" json:"organization_code"`
	CreditIdentificationCode     string `orm:"column(credit_identification_code)" form:"credit_identification_code" json:"credit_identification_code"`
	CompanyType                  string `orm:"column(company_type)" form:"company_type" json:"company_type"`
	TaxpayerIdentificationNumber string `orm:"column(taxpayer_identification_number)" form:"taxpayer_identification_number" json:"taxpayer_identification_number"`
	Trade                        string `orm:"column(trade)" form:"trade" json:"trade"`
	OperationPeriod              string `orm:"column(operation_period)" form:"operation_period" json:"operation_period"`
	RegistrationAuthority        string `orm:"column(registration_authority)" form:"registration_authority" json:"registration_authority"`
	RegisteredAddress            string `orm:"column(registered_address)" form:"registered_address" json:"registered_address"`
	ManagementScope              string `orm:"column(management_scope)" form:"management_scope" json:"management_scope"`
}

const (
	sqlForKeys          = `SELECT DISTINCT * FROM companies WHERE name LIKE '%%%s%%' OR legal_person LIKE '%%%s%%' OR business_registration_number LIKE '%%%s%%';`
	sqlForCreateCompany = "INSERT INTO companies(id,name,state,legal_person,registered_capital,registration_time,approval_time,business_registration_number,organization_code,credit_identification_code,company_type,taxpayer_identification_number,trade,operation_period,registration_authority,registered_address,management_scope) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	sqlForUpdateCompany = "UPDATE companies SET name=?, state=?, legal_person=?, registered_capital=?, registration_time=?, approval_time=?, business_registration_number=?, organization_code=?, credit_identification_code=?, company_type=?, taxpayer_identification_number=?, trade=?, operation_period=?, registration_authority=?, registered_address=?, management_scope=? WHERE id = ?"
)

func (c *Company) TableName() string {
	return "companies"
}

func getCompanyKey(id int) string {
	return fmt.Sprintf("/companies/%d", id)
}

func getCompaniesKey(keys []string) string {
	return strings.Join(keys, "")
}

func getCompanyFromRedis(id int) (*Company, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", getCompanyKey(id)))
	if err != nil {
		return nil, err
	}

	var com Company
	err = json.Unmarshal(b, &com)
	return &com, err
}

func setCompanyToRedis(c *Company) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", getCompanyKey(c.ID), b, "EX", EXTime)
	return err
}

func setCompaniesToRedis(keys []string, companies *[]Company) error {
	conn := RedisPool.Get()
	defer conn.Close()

	b, err := json.Marshal(companies)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", getCompaniesKey(keys), b, "EX", EXTime)
	return err
}

func getCompaniesFromReids(keys []string) (*[]Company, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	var companies []Company
	b, err := redis.Bytes(conn.Do("GET", getCompaniesKey(keys)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &companies)
	return &companies, err
}

func GetCompanyByID(id int) (*Company, error) {
	com, err := getCompanyFromRedis(id)
	if err != nil {
		var c Company
		err := O.Raw("SELECT * FROM companies WHERE id = ?", id).QueryRow(&c)
		if err == nil {
			setCompanyToRedis(&c)
		}
		return &c, err
	}

	return com, err
}

//添加，并去重
// func appendAndDistinct(src, dst *[]Company) {
// 	var x int
// 	for _, a := range *src {
// 		x = 0
// 		for _, b := range *dst {
// 			if reflect.DeepEqual(a, b) {
// 				x = 1
// 				break
// 			}
// 		}
// 		if x == 0 {
// 			*dst = append(*dst, a)
// 		}
// 	}
// }

func appendAndDistinct(coms1, coms2 *[]Company) []Company {
	var coms []Company
	for _, a := range *coms1 {
		for _, b := range *coms2 {
			if reflect.DeepEqual(a, b) {
				coms = append(coms, a)
				break
			}
		}
	}
	return coms
}

func GetCompaniesByKeys(keys []string) (*[]Company, error) {
	if len(keys) == 0 {
		return nil, errors.New("keys has zero key")
	}

	companies, err := getCompaniesFromReids(keys)
	if err != nil {
		var res []Company
		_, err := O.Raw(fmt.Sprintf(sqlForKeys, keys[0], keys[0], keys[0])).QueryRows(&res)
		if err != nil {
			return nil, err
		}

		var cs []Company
		for _, key := range keys[1:] {
			_, err := O.Raw(fmt.Sprintf(sqlForKeys, key, key, key)).QueryRows(&cs)
			if err != nil {
				return nil, err
			}
			res = appendAndDistinct(&cs, &res)
		}

		setCompaniesToRedis(keys, &res)
		return &res, nil
	}

	return companies, err
}

func CreateCompany(c *Company) (string, error) {
	if c == nil {
		return "", errors.New("arg c is nil")
	}
	var maxID int
	err := O.Raw("SELECT MAX(id) FROM companies;").QueryRow(&maxID)
	if err != nil {
		return "", err
	}

	_, err = O.Raw(sqlForCreateCompany,
		maxID+1,
		c.Name,
		c.State,
		c.LegalPerson,
		c.RegisteredCapital,
		c.RegistrationTime,
		c.ApprovalTime,
		c.BusinessRegistrationNumber,
		c.OrganizationCode,
		c.CreditIdentificationCode,
		c.CompanyType,
		c.TaxpayerIdentificationNumber,
		c.Trade,
		c.OperationPeriod,
		c.RegistrationAuthority,
		c.RegisteredAddress,
		c.ManagementScope).Exec()

	return fmt.Sprintf("%d", maxID), err
}

func defaultString(newStr, defaultStr string) string {
	if len(newStr) == 0 {
		return defaultStr
	}

	return newStr
}

func UpdateCompany(c *Company) error {

	oc, err := GetCompanyByID(c.ID)
	if err != nil {
		return err
	}

	c.Name = defaultString(c.Name, oc.Name)
	c.State = defaultString(c.State, oc.State)
	c.LegalPerson = defaultString(c.LegalPerson, oc.LegalPerson)
	c.RegisteredCapital = defaultString(c.RegisteredCapital, oc.RegisteredCapital)
	c.RegistrationTime = defaultString(c.RegistrationTime, oc.RegistrationTime)
	c.ApprovalTime = defaultString(c.ApprovalTime, oc.ApprovalTime)
	c.BusinessRegistrationNumber = defaultString(c.BusinessRegistrationNumber, oc.BusinessRegistrationNumber)
	c.OrganizationCode = defaultString(c.OrganizationCode, oc.OrganizationCode)
	c.CreditIdentificationCode = defaultString(c.CreditIdentificationCode, oc.CreditIdentificationCode)
	c.CompanyType = defaultString(c.CompanyType, oc.CompanyType)
	c.TaxpayerIdentificationNumber = defaultString(c.TaxpayerIdentificationNumber, oc.TaxpayerIdentificationNumber)
	c.Trade = defaultString(c.Trade, oc.Trade)
	c.OperationPeriod = defaultString(c.OperationPeriod, oc.OperationPeriod)
	c.RegistrationAuthority = defaultString(c.RegistrationAuthority, oc.RegistrationAuthority)
	c.RegisteredAddress = defaultString(c.RegisteredAddress, oc.RegisteredAddress)
	c.ManagementScope = defaultString(c.ManagementScope, oc.ManagementScope)

	_, err = O.Raw(sqlForUpdateCompany,
		c.Name,
		c.State,
		c.LegalPerson,
		c.RegisteredCapital,
		c.RegistrationTime,
		c.ApprovalTime,
		c.BusinessRegistrationNumber,
		c.OrganizationCode,
		c.CreditIdentificationCode,
		c.CompanyType,
		c.TaxpayerIdentificationNumber,
		c.Trade,
		c.OperationPeriod,
		c.RegistrationAuthority,
		c.RegisteredAddress,
		c.ManagementScope,
		c.ID).Exec()
	return err
}

func DeleteCompany(id int) error {
	_, err := O.Raw("DELETE FROM companies WHERE id = ?;", id).Exec()
	return err
}
