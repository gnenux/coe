package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Company struct {
	ID                           int    `orm:"column(id)"`
	Name                         string `orm:"column(name)"`
	State                        string `orm:"column(state)"`
	Points                       int    `orm:"column(points)"`
	LegalPerson                  string `orm:"column(legal_person)"`
	RegisteredCapital            string `orm:"column(registered_capital)"`
	RegistrationTime             string `orm:"column(registration_time)"`
	ApprovalTime                 string `orm:"column(approval_time)"`
	BusinessRegistrationNumber   string `orm:"column(business_registration_number)"`
	OrganizationCode             string `orm:"column(organization_code)"`
	CreditIdentificationCode     string `orm:"column(credit_identification_code)"`
	CompanyType                  string `orm:"column(company_type)"`
	TaxpayerIdentificationNumber string `orm:"column(taxpayer_identification_number)"`
	Trade                        string `orm:"column(trade)"`
	OperationPeriod              string `orm:"column(operation_period)"`
	RegistrationAuthority        string `orm:"column(registration_authority)"`
	RegisteredAddress            string `orm:"column(registered_address)"`
	ManagementScope              string `orm:"column(management_scope)"`
}

func (c *Company) TableName() string {
	return "companies"
}

const (
	sqlForKeys   = `SELECT DISTINCT * FROM companies WHERE name LIKE '%%%s%%' OR legal_person LIKE '%%%s%%' OR business_registration_number LIKE '%%%s%%';`
	sqlForCreate = "INSERT INTO companies(id,name,state,legal_person,registered_capital,registration_time,approval_time,business_registration_number,organization_code,credit_identification_code,company_type,taxpayer_identification_number,trade,operation_period,registration_authority,registered_address,management_scope) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	sqlForUpdate = "UPDATE companies SET name=?, state=?, legal_person=?, registered_capital=?, registration_time=?, approval_time=?, business_registration_number=?, organization_code=?, credit_identification_code=?, company_type=?, taxpayer_identification_number=?, trade=?, operation_period=?, registration_authority=?, registered_address=?, management_scope=? WHERE id = ?"
)

var O orm.Ormer

func init() {
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/tianyancha?charset=utf8")
	O = orm.NewOrm()
}

func GetCompanyById(id int) (*Company, error) {
	var c Company
	err := O.Raw("SELECT * FROM companies WHERE id = ?", id).QueryRow(&c)
	return &c, err
}

//添加，并去重
func appendAndDistinct(src, dst *[]Company) {
	var x int
	for _, a := range *src {
		x = 0
		for _, b := range *dst {
			if reflect.DeepEqual(a, b) {
				x = 1
				break
			}
		}
		if x == 0 {
			*dst = append(*dst, a)
		}
	}
}

func GetCompaniesByKeys(keys []string) (*[]Company, error) {
	var res []Company
	var cs []Company
	for _, key := range keys {
		_, err := O.Raw(fmt.Sprintf(sqlForKeys, key, key, key)).QueryRows(&cs)
		if err != nil {
			return nil, err
		}
		appendAndDistinct(&cs, &res)
	}
	return &res, nil
}

func CreateCompany(c *Company) error {
	if c == nil {
		return errors.New("arg c is nil")
	}
	var maxID int
	err := O.Raw("SELECT MAX(id) FROM companies;").QueryRow(&maxID)
	if err != nil {
		return err
	}

	_, err = O.Raw(sqlForCreate,
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

	return err
}

func UpdateCompany(c *Company) error {
	_, err := O.Raw(sqlForUpdate,
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

func main() {
	c, err := GetCompanyById(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *c)

	// res, err := GetCompaniesByKeys([]string{"上海", "北京"})
	// if err != nil {
	// 	panic(err)
	// }

	// for _, v := range *res {
	// 	fmt.Println(v.Name)
	// }
	c.Name = "博鳌财富（北京）国际文化传播有限公司"
	err = UpdateCompany(c)
	if err != nil {
		panic(err)
	}
	c, err = GetCompanyById(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *c)

}
