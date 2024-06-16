package models

// CountryCurrency [...]
type CountryCurrency struct {
	ID           int64    `gorm:"primaryKey;column:id" json:"-"`
	CountryCode  string   `gorm:"column:country_code" json:"countryCode"`
	Country      Country  `gorm:"joinForeignKey:country_code;foreignKey:country_code;references:CountryCode" json:"countryList"`
	CurrencyCode string   `gorm:"column:currency_code" json:"currencyCode"`
	Currency     Currency `gorm:"joinForeignKey:currency_code;foreignKey:currency_code;references:CurrencyCode" json:"currencyList"`
}

// TableName get sql table name.
func (m *CountryCurrency) TableName() string {
	return "country_currency"
}

// CountryCurrencyColumns get sql column name.
var CountryCurrencyColumns = struct {
	ID           string
	CountryCode  string
	CurrencyCode string
}{
	ID:           "id",
	CountryCode:  "country_code",
	CurrencyCode: "currency_code",
}
