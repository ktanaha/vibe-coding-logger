package writer

import (
	"fmt"
	"os"
	"sync"
	"vibe-coding-logger/pkg/logger"
	"vibe-coding-logger/internal/formatter"
)

// エイリアスを定義
type Entry = logger.Entry
type LogLevel = logger.LogLevel
type Formatter = logger.Formatter

// 定数のエイリアス
const (
	DEBUG = logger.DEBUG
	INFO  = logger.INFO
	WARN  = logger.WARN
	ERROR = logger.ERROR
	FATAL = logger.FATAL
)

// ConsoleWriter はコンソールへのログ出力を行う
type ConsoleWriter struct {
	formatter Formatter
	mu        sync.Mutex
}

// NewConsoleWriter は新しいコンソールライターを作成する
func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// Write はエントリをコンソールに書き込む
func (w *ConsoleWriter) Write(entry *Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.formatter == nil {
		// デフォルトのフォーマッターを使用
		w.formatter = formatter.NewTextFormatter()
	}

	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	// エラーレベルによって出力先を変える
	if entry.Level >= ERROR {
		_, err = os.Stderr.Write(formatted)
	} else {
		_, err = os.Stdout.Write(formatted)
	}

	return err
}

// SetFormatter はフォーマッターを設定する
func (w *ConsoleWriter) SetFormatter(f Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = f
}

// Close はライターを閉じる
func (w *ConsoleWriter) Close() error {
	return nil
}