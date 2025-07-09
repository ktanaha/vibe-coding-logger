package writer

import (
	"os"
	"sync"
	"vibe-coding-logger/internal"
	"vibe-coding-logger/internal/formatter"
)

// ConsoleWriter はコンソールへのログ出力を行う
type ConsoleWriter struct {
	formatter internal.Formatter
	mu        sync.Mutex
}

// NewConsoleWriter は新しいコンソールライターを作成する
func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// Write はエントリをコンソールに書き込む
func (w *ConsoleWriter) Write(entry *internal.Entry) error {
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
	if entry.Level >= internal.ERROR {
		_, err = os.Stderr.Write(formatted)
	} else {
		_, err = os.Stdout.Write(formatted)
	}

	return err
}

// SetFormatter はフォーマッターを設定する
func (w *ConsoleWriter) SetFormatter(f internal.Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = f
}

// Close はライターを閉じる
func (w *ConsoleWriter) Close() error {
	return nil
}
