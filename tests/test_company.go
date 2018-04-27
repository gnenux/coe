package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	b, err := json.Marshal(&Company{})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
