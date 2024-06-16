package database

import (
	"fmt"
	gormLogger "gorm.io/gorm/logger"
	"strings"
	"sync"

	// Configs
	cfg "github.com/nishanthrk/aspire-lms/app/configs"

	"github.com/nishanthrk/aspire-lms/app/logger"
	// Gorm
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// MysqlDB is the mysql connection handle
	MysqlDB   *gorm.DB
	onceMysql sync.Once
)

func ConnectMysql() {
	onceMysql.Do(func() {
		dsn := cfg.GetConfig().Mysql.GetMysqlConnectionInfo()
		logger.Sugar.Info(dsn)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
		db.Set("gorm:auto_preload", true)
		if err != nil {
			logger.Sugar.Info(strings.Repeat("!", 40))
			logger.Sugar.Info("😏 Could Not Establish Mysql DB Connection")
			logger.Sugar.Info(strings.Repeat("!", 40))
			logger.Sugar.Fatal(err)
		}

		// Set maximum number of open connections
		/*sqlDB, err := db.DB()
		if err != nil {
			logger.Sugar.Info(strings.Repeat("!", 40))
			logger.Sugar.Info("😏 Set maximum number of open connections")
			logger.Sugar.Info(strings.Repeat("!", 40))
			log.Fatal(err)
		}
		sqlDB.SetMaxOpenConns(10)

		// Set maximum number of idle connections
		sqlDB.SetMaxIdleConns(5)*/

		logger.Sugar.Info(strings.Repeat("-", 40))
		logger.Sugar.Info("😀 Connected To Mysql DB")
		logger.Sugar.Info(strings.Repeat("-", 40))

		MysqlDB = db
	})
}

type WhereCondition struct {
	Key            string               `json:"key"`
	Condition      string               `json:"condition"`
	Value          interface{}          `json:"value"`
	GroupCondition *GroupWhereCondition `json:"group_condition"`
	SubQuery       *SubQueryCondition   `json:"sub_query"`
}

type GroupWhereCondition struct {
	Condition string           `json:"condition"`
	Values    []WhereCondition `json:"values"`
}

type SubQueryCondition struct {
	TableName  string           `json:"table_name"`
	Model      interface{}      `json:"models"`
	FieldName  string           `json:"field_name"`
	Conditions []WhereCondition `json:"conditions"`
}

func ConditionBuilder(db *gorm.DB, whereConditions *[]WhereCondition, orConditions *[]WhereCondition, accessConditions *[]WhereCondition) *gorm.DB {
	if whereConditions != nil {
		for _, condition := range *whereConditions {
			if condition.GroupCondition != nil {
				var groupWhereClause string
				var groupArgs []interface{}
				for _, subCondition := range condition.GroupCondition.Values {
					if strings.HasPrefix(subCondition.Condition, "IS") {
						groupWhereClause += fmt.Sprintf("%s %s %s %s ", subCondition.Key, subCondition.Condition, subCondition.Value, condition.GroupCondition.Condition)
					} else {
						groupWhereClause += fmt.Sprintf("%s %s ? %s ", subCondition.Key, subCondition.Condition, condition.GroupCondition.Condition)
						groupArgs = append(groupArgs, subCondition.Value)
					}
				}

				if len(groupWhereClause) > 0 {
					groupWhereClause = strings.TrimSuffix(groupWhereClause, fmt.Sprintf(" %s ", condition.GroupCondition.Condition))
					db = db.Where(groupWhereClause, groupArgs...)
				}
			} else if condition.SubQuery != nil {
				// Process subQuery recursively
				subQuery := buildSubQuery(condition.SubQuery)
				// Modify the main query based on the subQuery result
				db = db.Where(fmt.Sprintf("%s %s (%s)", condition.Key, condition.Condition, subQuery))
			} else {
				if strings.HasPrefix(condition.Condition, "IS") {
					db = db.Where(fmt.Sprintf("%s %s %s", condition.Key, condition.Condition, condition.Value))
				} else {
					db = db.Where(fmt.Sprintf("%s %s ?", condition.Key, condition.Condition), condition.Value)
				}
			}
		}
	}

	if orConditions != nil {
		var orWhereClause string
		var orArgs []interface{}
		for _, condition := range *orConditions {
			if condition.SubQuery != nil {
				orWhereClause += fmt.Sprintf("%s %s (%s) OR ",
					condition.Key, condition.Condition, buildSubQuery(condition.SubQuery))
			} else if strings.HasPrefix(condition.Condition, "IS") {
				orWhereClause += fmt.Sprintf("%s %s %s OR ", condition.Key, condition.Condition, condition.Value)
			} else {
				orWhereClause += fmt.Sprintf("%s %s ? OR ", condition.Key, condition.Condition)
				orArgs = append(orArgs, condition.Value)
			}
		}

		// Remove the last " AND " or " OR " from orWhereClause
		if len(orWhereClause) > 0 {
			orWhereClause = orWhereClause[:len(orWhereClause)-4]
			db = db.Where(orWhereClause, orArgs...)
		}
	}

	if accessConditions != nil {
		var accessWhereClause string
		var accArgs []interface{}
		for _, condition := range *accessConditions {
			if condition.SubQuery != nil {
				subQuery := buildSubQuery(condition.SubQuery)
				accessWhereClause += fmt.Sprintf("%s %s (%s) OR ", condition.Key, condition.Condition, subQuery)
			} else if strings.HasPrefix(condition.Condition, "IS") {
				accessWhereClause += fmt.Sprintf("%s %s %s OR ", condition.Key, condition.Condition, condition.Value)
			} else {
				accessWhereClause += fmt.Sprintf("%s %s ? OR ", condition.Key, condition.Condition)
				accArgs = append(accArgs, condition.Value)
			}
		}

		// Remove the last " AND " or " OR " from accessWhereClause
		if len(accessWhereClause) > 0 {
			accessWhereClause = accessWhereClause[:len(accessWhereClause)-4]
			db = db.Where(accessWhereClause, accArgs...)
		}
	}

	return db
}

func buildSubQuery(subQuery *SubQueryCondition) string {
	sql := MysqlDB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = tx.Table(subQuery.TableName).Select(subQuery.FieldName)
		tx = ConditionBuilder(tx, &subQuery.Conditions, nil, nil)

		return tx.Find(subQuery.Model)
	})
	return sql
}
