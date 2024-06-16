package models

import "time"

// RepaymentPaymentLog [...]
type RepaymentPaymentLog struct {
	LogID       string    `gorm:"primaryKey;column:log_id" json:"-"`
	RepaymentID string    `gorm:"column:repayment_id" json:"repaymentId"`
	Repayment   Repayment `gorm:"joinForeignKey:repayment_id;foreignKey:repayment_id;references:RepaymentID" json:"repaymentList"`
	PaymentID   string    `gorm:"column:payment_id" json:"paymentId"`
	Payment     Payment   `gorm:"joinForeignKey:payment_id;foreignKey:payment_id;references:PaymentID" json:"paymentList"`
	Amount      float64   `gorm:"column:amount" json:"amount"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *RepaymentPaymentLog) TableName() string {
	return "repayment_payment_log"
}

// RepaymentPaymentLogColumns get sql column name.
var RepaymentPaymentLogColumns = struct {
	LogID       string
	RepaymentID string
	PaymentID   string
	Amount      string
	CreatedAt   string
	UpdatedAt   string
}{
	LogID:       "log_id",
	RepaymentID: "repayment_id",
	PaymentID:   "payment_id",
	Amount:      "amount",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}
