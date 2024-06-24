package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xenochrony/xylitol"
	"github.com/xenochrony/xylitol/logger"
	"github.com/xenochrony/xylitol/middleware"
	"github.com/xenochrony/xylitol/response"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestXylitolLogger(t *testing.T) {
	xylitol.Default.Debug("This is a log test for DEBUG level")
	xylitol.Default.Info("This is a log test for INFO level")
	xylitol.Default.Warn("This is a log test for WARN level")
	xylitol.Default.Error("This is a log test for ERROR level")
	// Fatal 会导致程序退出，因此不做测试
	//xylitol.Default.Fatal("This is a log test for FATAL level")
}

func TestXylitolLoggerWithFields(t *testing.T) {

	var fields []zapcore.Field

	xylitol.Default.Debug("This is a log test for DEBUG level", append(fields, zap.String("DEBUG", "debug"))...)
	xylitol.Default.Info("This is a log test for INFO level", append(fields, zap.String("INFO", "info"))...)
	xylitol.Default.Warn("This is a log test for WARN level", append(fields, zap.String("WARN", "warn"))...)
	xylitol.Default.Error("This is a log test for ERROR level", append(fields, zap.String("ERROR", "error"))...)
	// Fatal 会导致程序退出，因此不做测试
	//xylitol.Default.Fatal("This is a log test for FATAL level", append(fields, zap.String("FATAL", "fatal"))...)
}

func TestXylitolCustomLogger(t *testing.T) {
	logLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

	// TODO： 后期可以使用级别切换函数来动态切换日志级别
	for _, logLevel := range logLevels {
		logfile := path.Join("tmp", logLevel+".log")
		var Logger logger.LogInterface = xylitol.New(logLevel, logfile)
		xylitol.Default = Logger.(*logger.GLogger)
		fmt.Print(logLevel + "-----------------------------------------------\n")
		xylitol.Default.Debug("This is a log test for DEBUG level")
		xylitol.Default.Info("This is a log test for INFO level")
		xylitol.Default.Warn("This is a log test for WARN level")
		xylitol.Default.Error("This is a log test for ERROR level")
		// 注意：Fatal 会结束程序，根据需要决定是否取消注释
		// xylitol.Default.Fatal("This is a log test for FATAL level")
	}
}

func TestXylitolWithMiddleware(t *testing.T) {
	xylitol.Default = xylitol.New("debug", "tmp/run.log")

	testPath := "/ping"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine := gin.New()
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.GET(testPath, func(c *gin.Context) {
		c.JSON(204, nil)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	engine.ServeHTTP(res1, req1)
}

func TestXylitolWithMiddlewareAndResponse(t *testing.T) {
	xylitol.Default = xylitol.New("debug", "tmp/run.log")

	testPath := "/ping"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine := gin.New()
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	engine.GET(testPath, func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok", "data": "Hello"}, "Successfully logged in")
		response.Error(c, http.StatusInternalServerError, fmt.Errorf("hello"), "Failed logged in")
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	engine.ServeHTTP(res1, req1)
}

func TestXylitolWithGORM(t *testing.T) {
	xylitol.Default = xylitol.New("debug", "tmp/run.log")

	// 测试报错
	db, err := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"), middleware.GetGormConfig("debug"))

	// 测试不报错
	db, err = gorm.Open(sqlite.Open("gorm.db"), middleware.GetGormConfig("error"))
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}
	if err != nil {
		xylitol.Default.Error(err.Error())
	}
	// 迁移 schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		xylitol.Default.Fatal(err.Error())
	}

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product

	db.First(&product, "code = ?", "D43")   // 查找不存在的值
	db.Model(&product).Update("Price", 200) // 对不存在的记录进行更新
	// 更新多个不存在字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// Delete - 删除 product
	db.Delete(&product, 1)

	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)

}
