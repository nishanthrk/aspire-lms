package models

import (
	"github.com/nishanthrk/aspire-lms/app/common/constants"
	"github.com/nishanthrk/aspire-lms/app/common/utility"
	"github.com/nishanthrk/aspire-lms/app/database"
	"gorm.io/gorm"
	"time"
)

// User [...]
type User struct {
	UserID       string    `gorm:"primaryKey;column:user_id" json:"-"`
	UserName     string    `gorm:"column:user_name" json:"userName"`
	UserEmail    string    `gorm:"column:user_email" json:"userEmail"`
	UserPassword string    `gorm:"column:user_password" json:"userPassword"`
	UserType     string    `gorm:"column:user_type" json:"userType"`
	MobileNumber string    `gorm:"column:mobile_number" json:"mobileNumber"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *User) TableName() string {
	return "user"
}

// UsersColumns get sql column name.
var UsersColumns = struct {
	UserID       string
	UserName     string
	UserEmail    string
	UserPassword string
	UserType     string
	MobileNumber string
	CreatedAt    string
	UpdatedAt    string
}{
	UserID:       "user_id",
	UserName:     "user_name",
	UserEmail:    "user_email",
	UserPassword: "user_password",
	UserType:     "user_type",
	MobileNumber: "mobile_number",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UserPassword == "" {
		m.UserPassword = utility.HashPassword("12345678")
	}
	return
}

func (m *User) FindByPrimaryKey(userId string) (result User, err error) {
	err = database.MysqlDB.Model(m).Where("user_id = ?", userId).Find(&result).Error
	return
}

func (m *User) FindOneByCondition(whereCondition *[]database.WhereCondition, orCondition *[]database.WhereCondition) (result User, err error) {
	db := database.MysqlDB.Model(m)
	db = database.ConditionBuilder(db, whereCondition, orCondition, nil)
	err = db.Find(&result).Error
	return
}

func (m *User) FindLeastLoadedEmployee() (result string, err error) {
	var leastLoadedUser struct {
		UserID string
		Count  int64
	}
	err = database.MysqlDB.Model(m).
		Select("user.user_id, COUNT(loan_application_participant.participant_id) as count").
		Joins("LEFT JOIN loan_application_participant ON loan_application_participant.user_id = "+
			"user.user_id AND loan_application_participant.participant_type = ?", constants.UserTypeEmployee).
		Where("user.user_type = ?", constants.UserTypeEmployee).
		Group("user.user_id").
		Order("count ASC").
		Limit(1).
		Scan(&leastLoadedUser).Error

	result = leastLoadedUser.UserID
	return
}
