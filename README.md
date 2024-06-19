# xLogger 模块

本文档综合介绍了`GLogger`、`GMiddleware`及`GResponse`三个包的整合应用，旨在提升Go项目的日志管理和HTTP响应处理能力，确保服务的健壮性和可观察性。

## GLogger模块

### 功能

- **日志级别管理**：通过`GetLogLevel`动态解析并设置日志级别，兼容字符串配置。
- **日志记录**：提供结构化日志记录功能，包括生成唯一`traceId`。
- **输出配置**：支持日志输出到控制台与文件，并通过`lumberjack`实现日志文件的自动切割。

### 使用示例

```go
package main

import "github.com/xInitialization/xBackstage/internal/pkg/xLogger/GLogger"

func main() {
	GLogger.Logger.Debug("Hello World")
	GLogger.Logger.Info("Hello World"
	GLogger.Logger.Warn("Hello World")
	GLogger.Logger.Error("Hello World")
	GLogger.Logger.Fatal("Hello World")
}
```

默认的输出级别为 `Debug`，可以通过`GLogger.SetLogLevel`函数动态调整日志级别。如果你想设置为`Info`级别，你可以使用如下代码：

````go
package bootstrap

import (
	"github.com/xInitialization/xBackstage/configs"
	"github.com/xInitialization/xBackstage/internal/pkg/xLogger/GLogger"
	"path"
)

var Logger GLogger.LogInterface

func InitializeLogger() {
	loglevel := configs.ServerConfig.Log.Level
	logfile := path.Join(configs.ServerConfig.Log.Dir, configs.ServerConfig.Log.Name)
	Logger = GLogger.New(loglevel, logfile)
	GLogger.Default = Logger.(*GLogger.GLogger)
}
````

这样你就可以使用自定义的日志记录器来记录日志。

## GMiddleware模块

### 功能

- **请求日志**：`GinLogger`中间件记录请求的基本信息，如请求方法、路径、客户端IP等。
- **恢复机制**：`GinRecovery`中间件优雅处理panic，记录错误日志并可选包含调用栈信息，确保服务稳定性。


## GResponse模块

### 功能

- **响应结构**：定义了一套响应结构体和接口，用于统一API响应格式。
- **成功/错误处理**：`OK`和`Error`函数分别处理成功和错误响应，自动设置响应码、消息及日志记录。


### 整合流程

1. **初始化**：使用`GLogger`初始化日志系统，配置日志级别、输出方式。
2. **中间件注册**：在Gin框架中注册`GinLogger`和`GinRecovery`中间件，实现请求日志记录和异常恢复。
3. **响应构建**：在业务逻辑中，利用`GResponse`封装响应，统一处理成功与错误情况，自动记录响应日志。

### 使用示例

```go
package Login

import (
	"fmt"
	"github.com/xInitialization/xBackstage/internal/pkg/Utils"
	"github.com/xInitialization/xBackstage/internal/pkg/xLogger/GResponse"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xInitialization/xBackstage/internal/app/bootstrap"
	"github.com/xInitialization/xBackstage/internal/app/common/sys/User"
	"github.com/xInitialization/xBackstage/internal/app/services"
)

func Login(c *gin.Context) {
	// 定义一个变量来接收解析后的请求体数据
	var loginReq LoginRequest

	// 使用ShouldBindJSON来读取并绑定JSON请求体到loginReq
	if err := c.ShouldBindJSON(&loginReq); err == nil {
		// 成功解析，现在可以访问loginReq.Username和loginReq.Password
		fmt.Printf("Received login request: %+v\n", loginReq)
	} else {
		// 解析失败，返回错误信息
		GResponse.Error(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	db := bootstrap.GetDBFromPool("mysql")
	var user User.User
	result := db.First(&user, "user_name = ? and password = ?", loginReq.Username, Utils.SHA512(loginReq.Password))
	if result.Error != nil {
		GResponse.Error(c, http.StatusUnauthorized, result.Error, "Invalid username or password")
		return
	}
	token, err := services.GenerateJWT(user)
	if err != nil {
		GResponse.Error(c, http.StatusInternalServerError, err, "Failed to generate token")
		return
	}
	GResponse.OK(c, gin.H{"status": "ok", "type": loginReq.Type, "currentAuthority": "admin", "token": token}, "Successfully logged in")
}
```

此处日志记录会使用`GLogger.Default`来记录日志，因此你需要在初始化时设置默认日志记录器`GLogger.Default`为你自定义的日志记录器。

通过以上模块的整合应用，项目能够实现全面的日志记录、异常安全处理及标准化的API响应，提升开发效率和系统运维的便利性。

## 待优化事项

1. **GLogger 包优化**
   - **日志级别预解析**：在应用启动时解析日志级别，避免每次日志记录时的重复错误检查。
   - **日志生成器单例化**：实施单例模式管理`GLogger`实例，减少内存占用和资源消耗。
   - **日志级别配置**：提供日志级别动态配置，支持配置文件或环境变量驱动。

2. **GMiddleware 包优化**
   - **性能增强**：为请求参数记录添加条件判断，仅在调试模式下执行，减轻生产环境负担。
   - **中间件简化**：设计组合中间件方法，整合`GinLogger`和`GinRecovery`，简化应用配置。

3. **GResponse 包优化**
   - **深拷贝使用**：在`Clone`方法中采用深拷贝，避免数据篡改风险，尤其是处理复杂数据结构时。
   - **错误信息精简**：优化错误处理逻辑，使`Msg`字段传达用户友好信息，而`Info`字段提供详细错误堆栈。

4. **通用建议**
   - **依赖管理**：定期维护依赖库，确保最新版本的性能和安全性。
   - **单元测试**：全面覆盖各模块的单元测试，确保软件质量。
   - **文档完善**：加强外部文档编写，明确配置与使用指南，提升开发者体验。
   - **错误处理细化**：在错误恢复逻辑中引入更精细的错误分类处理，优化日志记录策略。
