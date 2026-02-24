package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	tests := []struct {
		name           string
		signal         os.Signal
		requestDelay   time.Duration
		shutdownDelay  time.Duration
		expectComplete bool
	}{
		{
			name:           "SIGINT - 请求在超时前完成",
			signal:         syscall.SIGINT,
			requestDelay:   2 * time.Second,
			shutdownDelay:  100 * time.Millisecond,
			expectComplete: true,
		},
		{
			name:           "SIGTERM - 请求在超时前完成",
			signal:         syscall.SIGTERM,
			requestDelay:   2 * time.Second,
			shutdownDelay:  100 * time.Millisecond,
			expectComplete: true,
		},
		{
			name:           "SIGQUIT - 请求在超时前完成",
			signal:         syscall.SIGQUIT,
			requestDelay:   2 * time.Second,
			shutdownDelay:  100 * time.Millisecond,
			expectComplete: true,
		},
		{
			name:           "请求超时 - 强制退出",
			signal:         syscall.SIGINT,
			requestDelay:   15 * time.Second,
			shutdownDelay:  100 * time.Millisecond,
			expectComplete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试服务器
			mux := http.NewServeMux()
			requestCompleted := make(chan bool, 1)

			mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.requestDelay)
				w.WriteHeader(http.StatusOK)
				requestCompleted <- true
			})

			srv := &http.Server{
				Addr:    ":18080",
				Handler: mux,
			}

			// 启动服务器
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					t.Logf("服务器错误: %v", err)
				}
			}()

			// 等待服务器启动
			time.Sleep(100 * time.Millisecond)

			// 发送测试请求
			go func() {
				_, _ = http.Get("http://localhost:18080/test")
			}()

			// 等待请求开始处理
			time.Sleep(tt.shutdownDelay)

			// 发送关闭信号
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			shutdownComplete := make(chan error, 1)
			go func() {
				shutdownComplete <- srv.Shutdown(shutdownCtx)
			}()

			// 检查结果
			select {
			case <-requestCompleted:
				if !tt.expectComplete {
					t.Error("期望请求被中断，但请求完成了")
				}
			case err := <-shutdownComplete:
				if tt.expectComplete {
					t.Error("期望请求完成，但服务器提前关闭")
				}
				if err != nil && err != context.DeadlineExceeded {
					t.Errorf("关闭错误: %v", err)
				}
			case <-time.After(12 * time.Second):
				t.Error("测试超时")
			}

			// 清理
			_ = srv.Close()
		})
	}
}

func TestSignalHandling(t *testing.T) {
	tests := []struct {
		name   string
		signal os.Signal
	}{
		{"处理 SIGINT", syscall.SIGINT},
		{"处理 SIGTERM", syscall.SIGTERM},
		{"处理 SIGQUIT", syscall.SIGQUIT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sigChan := make(chan os.Signal, 1)
			setupSignalHandler(sigChan)

			// 模拟发送信号
			sigChan <- tt.signal

			// 验证信号被接收
			select {
			case sig := <-sigChan:
				if sig != tt.signal {
					t.Errorf("期望信号 %v, 得到 %v", tt.signal, sig)
				}
			case <-time.After(1 * time.Second):
				t.Error("信号处理超时")
			}
		})
	}
}
