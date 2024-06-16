package models

import (
	"github.com/nishanthrk/aspire-lms/app/database"
)

// Country [...]
type Country struct {
	CountryCode string `gorm:"primaryKey;column:country_code" json:"-"`
	CountryName string `gorm:"column:country_name" json:"countryName"`
}

// TableName get sql table name.
func (m *Country) TableName() string {
	return "country"
}

// CountryColumns get sql column name.
var CountryColumns = struct {
	CountryCode string
	CountryName string
}{
	CountryCode: "country_code",
	CountryName: "country_name",
}

func (m *Country) FindByPrimaryKey(primaryKey string) (result Country, err error) {
	err = database.MysqlDB.Model(m).Where("country_code = ?", primaryKey).First(&result).Error
	return
}

func (m *Country) FindOneByCondition(whereCondition []database.WhereCondition) (result Country, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&result).Error
	return
}

func (m *Country) FindAllByCondition(whereCondition []database.WhereCondition) (results []Country, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&results).Error
	return
}
