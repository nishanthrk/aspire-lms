package models

import (
	"github.com/nishanthrk/aspire-lms/app/database"
	"time"
)

// LoanApplicationParticipant [...]
type LoanApplicationParticipant struct {
	ParticipantID   string          `gorm:"primaryKey;column:participant_id" json:"-"`
	ApplicationID   string          `gorm:"column:application_id" json:"applicationId"`
	LoanApplication LoanApplication `gorm:"joinForeignKey:application_id;foreignKey:application_id;references:ApplicationID" json:"loanApplicationList"`
	ParticipantType string          `gorm:"column:participant_type" json:"participantType"`
	UserID          string          `gorm:"column:user_id" json:"userId"`
	Users           User            `gorm:"joinForeignKey:user_id;foreignKey:user_id;references:UserID" json:"usersList"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName get sql table name.
func (m *LoanApplicationParticipant) TableName() string {
	return "loan_application_participant"
}

// LoanApplicationParticipantColumns get sql column name.
var LoanApplicationParticipantColumns = struct {
	ParticipantID   string
	ApplicationID   string
	ParticipantType string
	UserID          string
	CreatedAt       string
	UpdatedAt       string
}{
	ParticipantID:   "participant_id",
	ApplicationID:   "application_id",
	ParticipantType: "participant_type",
	UserID:          "user_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

func (m *LoanApplicationParticipant) FindAllByCondition(whereCondition []database.WhereCondition) (
	results []LoanApplicationParticipant, err error) {
	db := database.MysqlDB.Model(m).Preload("LoanApplication")
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&results).Error
	return
}

func (m *LoanApplicationParticipant) FindOneByCondition(whereCondition []database.WhereCondition) (
	result LoanApplicationParticipant, err error) {
	db := database.MysqlDB.Model(m).Preload("LoanApplication")
	db = database.ConditionBuilder(db, &whereCondition, nil, nil)
	err = db.Find(&result).Error
	return
}
