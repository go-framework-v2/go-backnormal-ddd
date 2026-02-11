package test

import (
	"testing"
	// "github.com/go-framework-v2/go-backnormal-gen/gen"
)

// 遇见测试不通过的问题，可能是包版本问题，执行以下：
// go mod tidy
// go get gorm.io/gorm@v1.25.0
// go get gorm.io/plugin/dbresolver@v1.4.7
func TestGen(t *testing.T) {
	// // 仓储接口实现 gen代码, 使用自研生成工具
	// dsn := "liuhua:liuhua@2024@tcp(rm-bp1ec3mzh5md190s1do.mysql.rds.aliyuncs.com:3306)/business_db?charset=utf8mb4&parseTime=True&loc=Local"
	// tableList := []string{"biz_user", "biz_app", "biz_sms_code"}
	// dddDir := "/Users/huanlema/Documents/Code_codeup/be-business/go-backnormal-ddd/src/internal/infrastructure/persistence/mysql"
	// tablePrefix := ""
	// domainPath := "go-backnormal-ddd/src/internal/domain/identity"
	// modelPath := "go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
	// err := gen.GenDdd(dsn, tableList, dddDir, tablePrefix, domainPath, modelPath)
	// if err != nil {
	// 	t.Error(err)
	// }
}
