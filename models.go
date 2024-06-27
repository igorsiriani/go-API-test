package main

import (
	"time"

    "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)
 
type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
}
 
type Record struct {
    gorm.Model
    PartnerId string
	PartnerName string
	CustomerId string
	CustomerName string
	CustomerDomainName string
	CustomerCountry string
	MpnId int
	Tier2MpnId string
	InvoiceNumber string
	ProductId string
	SkuId string
	AvailabilityId string
	SkuName string
	ProductName string
	PublisherName string
	PublisherId string
	SubscriptionDescription string
	SubscriptionId string
	ChargeStartDate time.Time
	ChargeEndDate time.Time
	UsageDate time.Time
	MeterType string
	MeterCategory string
	MeterId string
	MeterSubCategory string
	MeterName string
	MeterRegion string
	Unit string
	ResourceLocation string
	ConsumedService string
	ResourceGroup string
	ResourceURI string
	ChargeType string
	UnitPrice float64
	Quantity float64
	UnitType string
	BillingPreTaxTotal float64
	BillingCurrency string
	PricingPreTaxTotal float64
	PricingCurrency string
	ServiceInfo1 string
	ServiceInfo2 string
	Tags string
	AdditionalInfo string
	EffectiveUnitPrice float64
	PCToBCExchangeRate int
	PCToBCExchangeRateDate time.Time
	EntitlementId string
	EntitlementDescription string
	PartnerEarnedCreditPercentage int
	CreditPercentage int
	CreditType string
	BenefitOrderId string
	BenefitId string
	BenefitType string
}

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
 
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

type SearchCriteria struct {
    Column string `json:"column"`
    Value1 string `json:"value1"`
    Value2 string `json:"value2,omitempty"`
}
 
type SearchRequest struct {
    Criteria []SearchCriteria `json:"criteria"`
}
