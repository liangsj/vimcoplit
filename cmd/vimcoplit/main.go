package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/liangsj/vimcoplit/internal/api"
	"github.com/liangsj/vimcoplit/internal/core"
)

func main() {
	// 解析命令行参数
	port := flag.Int("port", 8080, "服务器监听端口")
	flag.Parse()

	// 初始化核心服务
	coreService := core.NewService()

	// 初始化API处理器
	handler := api.NewHandler(coreService)

	// 设置HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("正在关闭服务器...")
		if err := server.Close(); err != nil {
			log.Printf("关闭服务器时出错: %v\n", err)
		}
	}()

	// 启动服务器
	log.Printf("VimCoplit 服务器启动在端口 %d\n", *port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("服务器错误: %v\n", err)
	}
}
