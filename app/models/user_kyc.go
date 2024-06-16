package models

import "time"

// UserKyc [...]
type UserKyc struct {
	KycID       string    `gorm:"primaryKey;column:kyc_id" json:"-"`
	UserID      string    `gorm:"column:user_id" json:"userId"`
	Users       User      `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"usersList"`
	KycType     string    `gorm:"column:kyc_type" json:"identificationType"`
	KycNumber   string    `gorm:"column:kyc_number" json:"identificationNumber"`
	CountryCode string    `gorm:"column:country_code" json:"countryCode"`
	Country     Country   `gorm:"joinForeignKey:country_code;foreignKey:country_code;references:CountryCode" json:"countryList"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *UserKyc) TableName() string {
	return "user_kyc"
}

// UserKycColumns get sql column name.
var UserKycColumns = struct {
	KycID       string
	UserID      string
	KycType     string
	KycNumber   string
	CountryCode string
	CreatedAt   string
	UpdatedAt   string
}{
	KycID:       "id",
	UserID:      "user_id",
	KycType:     "kyc_type",
	KycNumber:   "kyc_number",
	CountryCode: "country_code",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}
