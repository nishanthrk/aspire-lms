package models

import (
	"github.com/Rhymond/go-money"
	"github.com/guregu/null"
	"github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"time"
)

// Repayment [...]
type Repayment struct {
	RepaymentID        string          `gorm:"primaryKey;column:repayment_id" json:"-"`
	ApplicationID      string          `gorm:"column:application_id" json:"applicationId"`
	LoanApplication    LoanApplication `gorm:"joinForeignKey:application_id;foreignKey:application_id;references:ApplicationID" json:"-"`
	InstallmentNumber  int             `gorm:"column:installment_number" json:"installmentNumber"`
	PrincipleAmount    float64         `gorm:"column:principle_amount" json:"principleAmount"`
	InterestAmount     float64         `gorm:"column:interest_amount" json:"interestAmount"`
	InstallmentDate    time.Time       `gorm:"column:installment_date" json:"dueDate"`
	PaymentDate        null.Time       `gorm:"column:payment_date" json:"paymentDate"`
	AmountDue          float64         `gorm:"column:amount_due" json:"amountDue"`
	AmountPaid         float64         `gorm:"column:amount_paid" json:"amountPaid"`
	OutstandingBalance null.Float      `gorm:"column:outstanding_balance" json:"outstandingBalance"`
	Status             string          `gorm:"column:status" json:"status"`
	PaymentReference   string          `gorm:"column:payment_reference" json:"paymentReference"`
	CreatedAt          time.Time       `gorm:"column:created_at" json:"-"`
	UpdatedAt          time.Time       `gorm:"column:updated_at" json:"-"`
}

// TableName get sql table name.
func (m *Repayment) TableName() string {
	return "repayment"
}

// RepaymentColumns get sql column name.
var RepaymentColumns = struct {
	RepaymentID        string
	ApplicationID      string
	InstallmentNumber  string
	PrincipleAmount    string
	InterestAmount     string
	DueDate            string
	PaymentDate        string
	AmountDue          string
	AmountPaid         string
	OutstandingBalance string
	Status             string
	PaymentReference   string
	CreatedAt          string
	UpdatedAt          string
}{
	RepaymentID:        "repayment_id",
	ApplicationID:      "application_id",
	InstallmentNumber:  "installment_number",
	PrincipleAmount:    "principle_amount",
	InterestAmount:     "interest_amount",
	DueDate:            "due_date",
	PaymentDate:        "payment_date",
	AmountDue:          "amount_due",
	AmountPaid:         "amount_paid",
	OutstandingBalance: "outstanding_balance",
	Status:             "status",
	PaymentReference:   "payment_reference",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

func (m *Repayment) GetRepaymentDTO() (object dto.Repayment) {
	object.InstallmentNumber = m.InstallmentNumber
	object.PaymentDate = m.InstallmentDate.Format("2006-01-02")
	object.CurrencyCode = m.LoanApplication.CurrencyCode
	object.Principle = money.NewFromFloat(m.PrincipleAmount, object.CurrencyCode).Display()
	object.Interest = money.NewFromFloat(m.InterestAmount, object.CurrencyCode).Display()
	object.Emi = money.NewFromFloat(m.AmountDue, object.CurrencyCode).Display()
	if m.AmountPaid > 0 {
		object.AmountPaid = null.StringFrom(money.NewFromFloat(m.AmountPaid, object.CurrencyCode).Display())
	}
	if m.OutstandingBalance.Float64 > 0 {
		object.OutstandingBalance = money.NewFromFloat(m.OutstandingBalance.Float64, object.CurrencyCode).Display()
	}
	object.RepaymentStatus = m.Status
	return
}

func (m *Repayment) FindAllByCondition(whereCondition []database.WhereCondition) (
	results []Repayment, err error) {
	db := database.MysqlDB.Model(m).Preload("LoanApplication")
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	db.Order("installment_date asc")
	err = db.Find(&results).Error
	return
}
