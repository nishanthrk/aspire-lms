package models

import "github.com/nishanthrk/aspire-lms/app/database"

// Currency [...]
type Currency struct {
	CurrencyCode string `gorm:"primaryKey;column:currency_code" json:"-"`
	CurrencyName string `gorm:"column:currency_name" json:"currencyName"`
}

// TableName get sql table name.
func (m *Currency) TableName() string {
	return "currency"
}

// CurrencyColumns get sql column name.
var CurrencyColumns = struct {
	CurrencyCode string
	CurrencyName string
}{
	CurrencyCode: "currency_code",
	CurrencyName: "currency_name",
}

func (m *Currency) FindByPrimaryKey(primaryKey string) (result Currency, err error) {
	err = database.MysqlDB.Model(m).Where("currency_code = ?", primaryKey).First(&result).Error
	return
}

func (m *Currency) FindOneByCondition(whereCondition []database.WhereCondition) (result Currency, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&result).Error
	return
}

func (m *Currency) FindAllByCondition(whereCondition []database.WhereCondition) (results []Currency, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&results).Error
	return
}
