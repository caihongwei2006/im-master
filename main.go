package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"im-master/router"
	"im-master/utils"

	// 导入生成的文档
	_ "im-master/docs"
)

// @title           IM System API
// @version         1.0
// @description     IM系统API文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.yourcompany.com/support
// @contact.email  support@yourcompany.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	// 记录开始时间
	startTime := time.Now()
	fmt.Printf("服务器启动于: %s\n", startTime.Format("2006-01-02 15:04:05"))

	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	if utils.DB == nil {
		panic("Failed to initialize database connection")
	}

	r := router.Router()

	// 创建通道来接收操作系统信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的 goroutine 中启动服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("监听错误: %s\n", err)
		}
	}()
	fmt.Println("服务器已启动，按 Ctrl+C 停止服务")

	// 阻塞，直到收到退出信号
	<-quit
	fmt.Println("正在关闭服务器...")

	// 设置10秒超时来优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("服务器强制关闭: %s\n", err)
	}

	// 计算服务器运行时间
	duration := time.Since(startTime)
	fmt.Printf("服务器运行了 %s\n", duration)
}
