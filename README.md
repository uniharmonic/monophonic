<img  align="right" src="https://avatars.githubusercontent.com/u/168158486?s=200&v=4" height="200" alt="logo"/>

[![monophonic](https://readme-typing-svg.demolab.com?font=Pixelify+Sans&size=64&pause=1000&center=false&vCenter=true&random=false&width=435&height=200&lines=:=>+monophonic+<=:)](https://github.com/uniharmonic/monophonic)

# Monophonic | 单声道

本文档综合介绍了`GLogger`、`GMiddleware`及`GResponse`三个包的整合应用，旨在提升
Go 项目的日志管理和 HTTP 响应处理能力，确保服务的健壮性和可观察性。

## `GLogger` 模块

`GLogger` 是下述所有中间件的基础模块，提供了日志记录、日志级别管理、日志输出配置
等功能。

### 功能

- **日志级别管理**：通过`GetLogLevel`动态解析并设置日志级别，兼容字符串配置。
- **日志记录**：提供结构化日志记录功能，包括生成唯一`traceId`。
- **输出配置**：支持日志输出到控制台与文件，并通过`lumberjack`实现日志文件的自动
  切割。

### 使用示例

#### 日志初始化

模块引入时会自动初始化默认日志记录器，默认的日志记录级别为 `Debug`， 如果需要自
定义日志记录器，可以在按照如下方式进行初始化。

```go
package bootstrap

import (
	"path"

	"github.com/uniharmonic/xerography/configs"
	"github.com/uniharmonic/monophonic"
	"github.com/uniharmonic/monophonic/logger"
)

var Logger logger.LogInterface

func InitializeLogger() {
	// 此处可以修改为你自己对应的日志文件配置，例如从环境变量或者配置文件进行读取
	loglevel := "Info"
	logfile := "./tmp/run.log"
	// 此处使用 Logger 来进行日志管理，实际上你仍然可以使用 monophonic.Default 来进行日志管理
	Logger = monophonic.New(loglevel, logfile)
	monophonic.Default = Logger.(*logger.GLogger)
}
```

#### 输出日志

如果您需要手动输出某些日志，您可以使用`monophonic.Default`来输出日志。

默认的输出级别有`Debug`、`Info`、`Warn`、`Error`和`Fatal`五个级别。

```go
package main

import "github.com/uniharmonic/monophonic"

func main() {
	monophonic.Default.Debug("This is a log test for DEBUG level")
	monophonic.Default.Info("This is a log test for INFO level")
	monophonic.Default.Warn("This is a log test for WARN level")
	monophonic.Default.Error("This is a log test for ERROR level")
	// Fatal 会导致程序退出
	monophonic.Default.Fatal("This is a log test for FATAL level")
}
```

> 注意：`Fatal`级别的日志会直接导致程序退出，请谨慎使用。

#### 动态切换日志级别

除了一开始进行日志级别的初始化外，您还可以通过`GLogger.SetLogLevel`函数动态调整日志级别。

```go
monophonic.Default.SetLogLevel("Info")	// Info, Warn, Error, Fatal, Debug 均可（不区分大小写）
```

## Middleware（中间件）

### Gin 中间件

`GinLogger`和`GinRecovery`是`Gin`框架的中间件，用于接管 `Gin` 框架默认的记录请求日志和恢复程序异常。同时我们提供了 `GReturn` 中间件，用于返回自定义的响应结构体（统一返回的结构体）。

#### 功能

- **请求日志**：`GinLogger`中间件记录请求的基本信息，如请求方法、路径、客户端 IP
  等。
- **恢复机制**：`GinRecovery`中间件优雅处理 panic，记录错误日志并可选包含调用栈
  信息，确保服务稳定性。
- **Greturn**：
  - **响应结构**：定义了一套响应结构体和接口，用于统一 API 响应格式。
  - **成功/错误处理**：`OK`和`Error`函数分别处理成功和错误响应，自动设置响应码、消
    息及日志记录。

#### 使用示例

1. **初始化**：首先您需要使用`GLogger`初始化日志系统。
2. **中间件注册**：在 Gin 框架中注册`GinLogger`和`GinRecovery`中间件，实现请求日
   志记录和异常恢复。
3. **响应构建**：在业务逻辑中，利用`GResponse`封装响应，统一处理成功与错误情况，
   自动记录响应日志。

```go
package xgin

import (
	"fmt"
	"github.com/uniharmonic/monophonic/middleware"
	"github.com/gin-gonic/gin"
)

// New 初始化一个新的 xgin 实例并返回它
func main() {
	// 日志初始化
	_HERE_IS_YOUR_LOGGER_INIT_CODE_HERE()

	r = gin.New()
	// 注册日志中间件
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 注册自定义响应中间件
	r.POST("/api/v1/user/login", func(c *gin.Context) {
		// 定义一个变量来接收解析后的请求体数据
		var loginReq LoginRequest

		// 使用ShouldBindJSON来读取并绑定JSON请求体到loginReq
		if err := c.ShouldBindJSON(&loginReq); err == nil {
			// 成功解析，现在可以访问loginReq.Username和loginReq.Password
			fmt.Printf("Received login request: %+v\n", loginReq)
		} else {
			// 解析失败，返回错误信息
			response.Error(c, http.StatusBadRequest, err, "Invalid request body")
			return
		}

		// 数据库连接池获取数据库连接，请把此处替换为你的连接池获取数据库连接的代码
		db := bootstrap.GetDBFromPool("mysql")
		var user User.User
		result := db.First(&user, "user_name = ? and password = ?", loginReq.Username, Utils.SHA512(loginReq.Password))

		// 统一返回
		if result.Error != nil {
			response.Error(c, http.StatusUnauthorized, result.Error, "Invalid username or password")
			return
		}
		token, err := services.GenerateJWT(user)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err, "Failed to generate token")
			return
		}
		response.OK(c, gin.H{"status": "ok", "type": loginReq.Type, "currentAuthority": "admin", "token": token}, "Successfully logged in")
	}
}
```

> 此处日志记录会使用`monophonic.Default`来记录日志，因此你需要在初始化时设置默认日志记录器`monophonic.Default`为你自定义的日志记录器。

### GORM 中间件

`GORM`中间件用于记录`GORM`操作的日志，包括`SQL`语句、执行时间、参数等。

此处执行比较简单，您只需要在构建数据库连接池时使用 `middleware.GetGormConfig` 方法注册`GORM`中间件即可。

```go
// 其中 error 是 GORM 的配置选项，用于控制错误日志的记录级别。
// Info, Warn, Error, Fatal, Debug 均可（不区分大小写）。
// 实质上是调用了 monophonic.Default.SetLogLevel 方法来设置日志级别。
db, err = gorm.Open(sqlite.Open("gorm.db"), middleware.GetGormConfig("error"))
```

> 因为此处实质上是调用了 monophonic.Default.SetLogLevel 方法来设置日志级别，因此需要在初始化时设置默认日志记录器`monophonic.Default`为你自定义的日志记录器。
> 同时这样做实际上并不是真的修改了 GORM 的日志记录器，而是只是修改了默认的日志记录器导致其输出不显示而已，因此对于性能会有一定影响。

## 待优化事项

- [ ] GLogger 包优化
	- [ ] 优化日志级别预解析
	- [ ] 单例模式管理`GLogger`实例，减少内存占用和资源消耗。
	- [ ] 采用 Skipper 模式，允许用户自定义跳过某些请求的日志记录。
- [ ] 中间件包优化
	- [ ] 为请求参数记录添加条件判断，仅在调试模式下执行，减轻生产环境负担。
	- [ ] 在 `Response` 对象中，在`Clone`方法中采用深拷贝，避免数据篡改风险，尤其是处理复杂数据结构时。
	- [ ] 优化错误处理逻辑，使`Msg`字段传达用户友好信息，而`Info`字段提供详细错误堆栈。
- [ ] 在错误恢复逻辑中引入更精细的错误分类处理，优化日志记录策略。
- [ ] 优化单元测试，确保每个功能模块的测试覆盖率达到100%。