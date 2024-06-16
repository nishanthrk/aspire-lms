package models

import (
	"github.com/guregu/null"
	"github.com/nishanthrk/aspire-lms/app/common/utility"
	"github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"time"
)

// LoanApplication [...]
type LoanApplication struct {
	ApplicationID      string      `gorm:"primaryKey;column:application_id" json:"-"`
	LoanAmount         float64     `gorm:"column:loan_amount" json:"loanAmount"`
	CurrencyCode       string      `gorm:"column:currency_code" json:"currencyCode"`
	Currency           Currency    `gorm:"joinForeignKey:currency_code;foreignKey:currency_code;references:CurrencyCode" json:"currencyList"`
	InterestRate       float64     `gorm:"column:interest_rate" json:"interestRate"`
	LoanTerm           int         `gorm:"column:loan_term" json:"loanTerm"`
	LoanTermUnit       string      `gorm:"column:loan_term_unit" json:"loanTermUnit"`
	ApplicationDate    time.Time   `gorm:"column:application_date" json:"applicationDate"`
	Status             string      `gorm:"column:status" json:"status"`
	ApprovedAmount     null.Float  `gorm:"column:approved_amount" json:"approvedAmount"`
	ApprovedDate       null.Time   `gorm:"column:approved_date" json:"approvedDate"`
	RejectionReason    null.String `gorm:"column:rejection_reason" json:"rejectionReason"`
	RepaymentStartDate null.Time   `gorm:"column:repayment_start_date" json:"repaymentStartDate"`
	Income             float64     `gorm:"column:income" json:"income"`
	CreditScore        int         `gorm:"column:credit_score" json:"creditScore"`
	ExistingDebts      float64     `gorm:"column:existing_debts" json:"existingDebts"`
	CountryCode        string      `gorm:"column:country_code" json:"countryCode"`
	Country            Country     `gorm:"joinForeignKey:country_code;foreignKey:country_code;references:CountryCode" json:"countryList"`
	EligibleLoanAmount float64     `gorm:"column:eligible_loan_amount" json:"eligibleLoanAmount"`
	CreatedAt          time.Time   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time   `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *LoanApplication) TableName() string {
	return "loan_application"
}

// LoanApplicationColumns get sql column name.
var LoanApplicationColumns = struct {
	ApplicationID      string
	LoanAmount         string
	CurrencyCode       string
	InterestRate       string
	LoanTerm           string
	LoanTermUnit       string
	ApplicationDate    string
	Status             string
	ApprovedAmount     string
	ApprovedDate       string
	RejectionReason    string
	RepaymentStartDate string
	Income             string
	CreditScore        string
	ExistingDebts      string
	CountryCode        string
	EligibleLoanAmount string
	CreatedAt          string
	UpdatedAt          string
}{
	ApplicationID:      "application_id",
	LoanAmount:         "loan_amount",
	CurrencyCode:       "currency_code",
	InterestRate:       "interest_rate",
	LoanTerm:           "loan_term",
	LoanTermUnit:       "loan_term_unit",
	ApplicationDate:    "application_date",
	Status:             "status",
	ApprovedAmount:     "approved_amount",
	ApprovedDate:       "approved_date",
	RejectionReason:    "rejection_reason",
	RepaymentStartDate: "repayment_start_date",
	Income:             "income",
	CreditScore:        "credit_score",
	ExistingDebts:      "existing_debts",
	CountryCode:        "country_code",
	EligibleLoanAmount: "eligible_loan_amount",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

func (m *LoanApplication) FindByPrimaryKey(applicationId string) (result LoanApplication, err error) {
	err = database.MysqlDB.Model(m).Where("application_id = ?", applicationId).First(&result).Error
	return
}

func (m *LoanApplication) GetLoanApplicationDTO() (object dto.LoanApplicationObject) {
	object.ApplicationId = utility.ToString(m.ApplicationID)
	object.ApplicationStatus = m.Status
	object.LoanAmount = m.LoanAmount
	object.CurrencyCode = m.CurrencyCode
	object.InterestRate = m.InterestRate
	object.LoanTerm = m.LoanTerm
	object.LoanTermUnit = m.LoanTermUnit
	object.Income = m.Income
	object.CreditScore = m.CreditScore
	object.ExistingDebts = m.ExistingDebts
	object.CountryCode = m.CountryCode
	return
}
