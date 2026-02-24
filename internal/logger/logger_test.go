package logger

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name       string
		logDir     string
		maxBackups int
		wantErr    bool
	}{
		{
			name:       "valid config",
			logDir:     "testlogs",
			maxBackups: 30,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.RemoveAll(tt.logDir)
			})

			err := Init(tt.logDir, tt.maxBackups)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if _, err := os.Stat(tt.logDir); os.IsNotExist(err) {
					t.Errorf("log directory not created")
				}
			}
		})
	}
}

func TestLogger(t *testing.T) {
	logDir := "testlogs"
	t.Cleanup(func() {
		os.RemoveAll(logDir)
	})

	err := Init(logDir, 30)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	tests := []struct {
		name    string
		logFunc func(string, ...zap.Field)
		level   zapcore.Level
		message string
	}{
		{
			name:    "info log",
			logFunc: Info,
			level:   zapcore.InfoLevel,
			message: "test info",
		},
		{
			name:    "error log",
			logFunc: Error,
			level:   zapcore.ErrorLevel,
			message: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(tt.message)
			Sync()

			var logFile string
			if tt.level == zapcore.ErrorLevel {
				logFile = filepath.Join(logDir, "error.log")
			} else {
				logFile = filepath.Join(logDir, "info.log")
			}

			if _, err := os.Stat(logFile); os.IsNotExist(err) {
				t.Errorf("log file %s not created", logFile)
			}
		})
	}
}

func TestSync(t *testing.T) {
	logDir := "testlogs"
	t.Cleanup(func() {
		os.RemoveAll(logDir)
	})

	err := Init(logDir, 30)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	Info("test sync")
	err = Sync()
	// Sync may fail on stdout in tests, ignore that error
	if err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
		t.Errorf("Sync() unexpected error = %v", err)
	}
}
