package models

import (
	"github.com/nishanthrk/aspire-lms/app/database"
	"time"
)

// LoanEligibilityConfig [...]
type LoanEligibilityConfig struct {
	ID                      int64     `gorm:"primaryKey;column:id" json:"-"`
	CountryCode             string    `gorm:"column:country_code" json:"countryCode"`
	Country                 Country   `gorm:"joinForeignKey:country_code;foreignKey:country_code;references:CountryCode" json:"countryList"`
	MinCreditScore          int       `gorm:"column:min_credit_score" json:"minCreditScore"`
	MaxCreditScore          int       `gorm:"column:max_credit_score" json:"maxCreditScore"`
	MaxFoir                 float64   `gorm:"column:max_foir" json:"maxFoir"`
	BaseLoanAmount          float64   `gorm:"column:base_loan_amount" json:"baseLoanAmount"`
	CreditScoreFactorHigh   float64   `gorm:"column:credit_score_factor_high" json:"creditScoreFactorHigh"`
	CreditScoreFactorMedium float64   `gorm:"column:credit_score_factor_medium" json:"creditScoreFactorMedium"`
	CreditScoreFactorLow    float64   `gorm:"column:credit_score_factor_low" json:"creditScoreFactorLow"`
	CreatedAt               time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt               time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *LoanEligibilityConfig) TableName() string {
	return "loan_eligibility_config"
}

// LoanEligibilityConfigColumns get sql column name.
var LoanEligibilityConfigColumns = struct {
	ID                      string
	CountryCode             string
	MinCreditScore          string
	MaxCreditScore          string
	MaxFoir                 string
	BaseLoanAmount          string
	CreditScoreFactorHigh   string
	CreditScoreFactorMedium string
	CreditScoreFactorLow    string
	CreatedAt               string
	UpdatedAt               string
}{
	ID:                      "id",
	CountryCode:             "country_code",
	MinCreditScore:          "min_credit_score",
	MaxCreditScore:          "max_credit_score",
	MaxFoir:                 "max_foir",
	BaseLoanAmount:          "base_loan_amount",
	CreditScoreFactorHigh:   "credit_score_factor_high",
	CreditScoreFactorMedium: "credit_score_factor_medium",
	CreditScoreFactorLow:    "credit_score_factor_low",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
}

func (m *LoanEligibilityConfig) FindOneByCondition(whereCondition []database.WhereCondition) (result LoanEligibilityConfig, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&result).Error
	return
}
