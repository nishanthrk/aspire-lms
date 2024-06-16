package models

import "time"

// Payment [...]
type Payment struct {
	PaymentID       string          `gorm:"primaryKey;column:payment_id" json:"-"`
	CurrencyCode    string          `gorm:"column:currency_code" json:"currencyCode"`
	Currency        Currency        `gorm:"joinForeignKey:currency_code;foreignKey:currency_code;references:CurrencyCode" json:"currencyList"`
	Amount          float64         `gorm:"column:amount" json:"amount"`
	ApplicationID   string          `gorm:"column:application_id" json:"applicationId"`
	LoanApplication LoanApplication `gorm:"joinForeignKey:application_id;foreignKey:application_id;references:ApplicationID" json:"loanApplicationList"`
	Status          string          `gorm:"column:status" json:"status"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *Payment) TableName() string {
	return "payment"
}

// PaymentColumns get sql column name.
var PaymentColumns = struct {
	PaymentID     string
	CurrencyCode  string
	Amount        string
	ApplicationID string
	Status        string
	CreatedAt     string
	UpdatedAt     string
}{
	PaymentID:     "payment_id",
	CurrencyCode:  "currency_code",
	Amount:        "amount",
	ApplicationID: "application_id",
	Status:        "status",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}
