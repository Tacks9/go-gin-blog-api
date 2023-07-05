package models

import (
	"fmt"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// 模型基类
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// 初始化数据库实例
func Setup() {
	var (
		err error
		// dbType, dbName, user, password, host, tablePrefix string
	)

	// // 读取配置
	// sec, err := setting.Cfg.GetSection("database")
	// if err != nil {
	// 	log.Fatal(2, "Fail to get section 'database': %v", err)
	// }

	// dbType = sec.Key("TYPE").String()
	// dbName = sec.Key("NAME").String()
	// user = sec.Key("USER").String()
	// password = sec.Key("PASSWORD").String()
	// host = sec.Key("HOST").String()
	// tablePrefix = sec.Key("TABLE_PREFIX").String()

	// 打开数据库配置
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	))

	if err != nil {
		log.Println(err)
	}

	// 获取默认表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	fmt.Println(setting.DatabaseSetting)

	// 设置表名为单数
	db.SingularTable(true)
	// 开启日志
	db.LogMode(true)

	// 设置回调函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	// 设置删除回调
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	// 最大连接数
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

// 关闭数据库
func CloseDB() {
	defer db.Close()
}

// 更新时间
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()

		// 判断有无 created_on 字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			// 如果为空进行设置
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		// 判断有无 modified_on 字段
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新修改时间
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	// 如果没有 update_column ，就会去设置 modified_on 字段
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除记录
func deleteCallback(scope *gorm.Scope) {
	if scope.HasError() {
		logging.Info("deleteCallback Error")
		return
	}

	// 尝试获取 delete_option
	var extraOption string
	if str, ok := scope.Get("gorm:delete_option"); ok {
		extraOption = fmt.Sprint(str)
	}

	// 获取我们约定的删除字段，，
	deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

	if !scope.Search.Unscoped && hasDeletedOnField {
		// 若存在则 UPDATE 软删除
		scope.Raw(fmt.Sprintf(
			"UPDATE %v SET %v=%v %v %v",
			// 当前引用的表名
			scope.QuotedTableName(),
			// 删除的字段
			scope.Quote(deletedOnField.DBName),
			// 添加值作为 SQL 的参数
			scope.AddToVars(time.Now().Unix()),
			// 返回组合好的条件 SQL
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			// 设置删除时间
			addExtraSpaceIfExist(extraOption),
		)).Exec()
	} else {
		// 若不存在则 DELETE 硬删除
		scope.Raw(fmt.Sprintf(
			"DELETE FROM %v%v%v",
			scope.QuotedTableName(),
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
			addExtraSpaceIfExist(extraOption),
		)).Exec()
	}
}

// 新增空格
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
