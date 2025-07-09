package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"vibe-coding-logger/internal"
	"vibe-coding-logger/internal/formatter"
)

// FileWriter はファイルへのログ出力を行う
type FileWriter struct {
	filename  string
	file      *os.File
	formatter internal.Formatter
	mu        sync.Mutex
}

// NewFileWriter は新しいファイルライターを作成する
func NewFileWriter(filename string) (*FileWriter, error) {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		filename: filename,
		file:     file,
	}, nil
}

// Write はエントリをファイルに書き込む
func (w *FileWriter) Write(entry *internal.Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.formatter == nil {
		// デフォルトのフォーマッターを使用
		w.formatter = formatter.NewJSONFormatter()
	}

	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = w.file.Write(formatted)
	if err != nil {
		return err
	}

	// 即座に書き込み（バッファリングを無効化）
	return w.file.Sync()
}

// SetFormatter はフォーマッターを設定する
func (w *FileWriter) SetFormatter(f internal.Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = f
}

// Close はファイルを閉じる
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// RotatingFileWriter はローテーションするファイルライター
type RotatingFileWriter struct {
	baseFilename string
	maxSize      int64
	maxFiles     int
	currentFile  *os.File
	currentSize  int64
	formatter    internal.Formatter
	mu           sync.Mutex
}

// NewRotatingFileWriter は新しいローテーションファイルライターを作成する
func NewRotatingFileWriter(baseFilename string, maxSize int64, maxFiles int) (*RotatingFileWriter, error) {
	dir := filepath.Dir(baseFilename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	rfw := &RotatingFileWriter{
		baseFilename: baseFilename,
		maxSize:      maxSize,
		maxFiles:     maxFiles,
	}

	if err := rfw.openCurrentFile(); err != nil {
		return nil, err
	}

	return rfw, nil
}

// Write はエントリをローテーションファイルに書き込む
func (w *RotatingFileWriter) Write(entry *internal.Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.formatter == nil {
		w.formatter = formatter.NewJSONFormatter()
	}

	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	// ファイルサイズをチェック
	if w.currentSize+int64(len(formatted)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return err
		}
	}

	n, err := w.currentFile.Write(formatted)
	if err != nil {
		return err
	}

	w.currentSize += int64(n)
	return w.currentFile.Sync()
}

// rotate はファイルをローテーションする
func (w *RotatingFileWriter) rotate() error {
	// 現在のファイルを閉じる
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	// 既存のファイルを番号付きにリネーム
	for i := w.maxFiles - 1; i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d", w.baseFilename, i)
		newName := fmt.Sprintf("%s.%d", w.baseFilename, i+1)

		if i == w.maxFiles-1 {
			// 最古のファイルを削除
			os.Remove(newName)
		}

		if _, err := os.Stat(oldName); err == nil {
			_ = os.Rename(oldName, newName) // エラーは無視（ログローテーション時のベストエフォート）
		}
	}

	// 現在のファイルを .1 にリネーム
	if _, err := os.Stat(w.baseFilename); err == nil {
		_ = os.Rename(w.baseFilename, w.baseFilename+".1") // エラーは無視（ログローテーション時のベストエフォート）
	}

	// 新しいファイルを作成
	return w.openCurrentFile()
}

// openCurrentFile は現在のファイルを開く
func (w *RotatingFileWriter) openCurrentFile() error {
	file, err := os.OpenFile(w.baseFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file

	// 現在のファイルサイズを取得
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	w.currentSize = stat.Size()

	return nil
}

// SetFormatter はフォーマッターを設定する
func (w *RotatingFileWriter) SetFormatter(formatter internal.Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = formatter
}

// Close はファイルを閉じる
func (w *RotatingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

// DailyRotatingFileWriter は日付ベースでローテーションするファイルライター
type DailyRotatingFileWriter struct {
	baseFilename string
	currentDate  string
	currentFile  *os.File
	formatter    internal.Formatter
	mu           sync.Mutex
}

// NewDailyRotatingFileWriter は新しい日次ローテーションファイルライターを作成する
func NewDailyRotatingFileWriter(baseFilename string) (*DailyRotatingFileWriter, error) {
	dir := filepath.Dir(baseFilename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	drfw := &DailyRotatingFileWriter{
		baseFilename: baseFilename,
	}

	if err := drfw.openCurrentFile(); err != nil {
		return nil, err
	}

	return drfw, nil
}

// Write はエントリを日次ローテーションファイルに書き込む
func (w *DailyRotatingFileWriter) Write(entry *internal.Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.formatter == nil {
		w.formatter = formatter.NewJSONFormatter()
	}

	// 日付が変わったかチェック
	currentDate := time.Now().Format("2006-01-02")
	if w.currentDate != currentDate {
		if err := w.rotate(); err != nil {
			return err
		}
	}

	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = w.currentFile.Write(formatted)
	if err != nil {
		return err
	}

	return w.currentFile.Sync()
}

// rotate はファイルをローテーションする
func (w *DailyRotatingFileWriter) rotate() error {
	// 現在のファイルを閉じる
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	// 新しいファイルを開く
	return w.openCurrentFile()
}

// openCurrentFile は現在のファイルを開く
func (w *DailyRotatingFileWriter) openCurrentFile() error {
	currentDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s.%s", w.baseFilename, currentDate)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	w.currentDate = currentDate
	return nil
}

// SetFormatter はフォーマッターを設定する
func (w *DailyRotatingFileWriter) SetFormatter(formatter internal.Formatter) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatter = formatter
}

// Close はファイルを閉じる
func (w *DailyRotatingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

// VibeFileWriter はバイブコーディング専用のファイルライター
type VibeFileWriter struct {
	*FileWriter
	sessionID     string
	problemDomain string
}

// NewVibeFileWriter は新しいバイブファイルライターを作成する
func NewVibeFileWriter(baseFilename, sessionID, problemDomain string) (*VibeFileWriter, error) {
	// セッションIDと問題ドメインを含むファイル名を生成
	filename := fmt.Sprintf("%s_%s_%s.log", baseFilename, sessionID, problemDomain)

	fileWriter, err := NewFileWriter(filename)
	if err != nil {
		return nil, err
	}

	vfw := &VibeFileWriter{
		FileWriter:    fileWriter,
		sessionID:     sessionID,
		problemDomain: problemDomain,
	}

	// バイブ専用フォーマッターを設定
	vfw.SetFormatter(formatter.NewVibeJSONFormatter())

	return vfw, nil
}

// BufferedFileWriter はバッファリングされたファイルライター
type BufferedFileWriter struct {
	*FileWriter
	buffer     [][]byte
	bufferSize int
	mu         sync.Mutex
}

// NewBufferedFileWriter は新しいバッファリングファイルライターを作成する
func NewBufferedFileWriter(filename string, bufferSize int) (*BufferedFileWriter, error) {
	fileWriter, err := NewFileWriter(filename)
	if err != nil {
		return nil, err
	}

	return &BufferedFileWriter{
		FileWriter: fileWriter,
		buffer:     make([][]byte, 0, bufferSize),
		bufferSize: bufferSize,
	}, nil
}

// Write はエントリをバッファに書き込む
func (w *BufferedFileWriter) Write(entry *internal.Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.formatter == nil {
		w.formatter = formatter.NewJSONFormatter()
	}

	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}

	w.buffer = append(w.buffer, formatted)

	// バッファが満杯になったら出力
	if len(w.buffer) >= w.bufferSize {
		return w.flush()
	}

	return nil
}

// Flush はバッファの内容をファイルに出力する
func (w *BufferedFileWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flush()
}

// flush はバッファの内容をファイルに出力する（内部用）
func (w *BufferedFileWriter) flush() error {
	if len(w.buffer) == 0 {
		return nil
	}

	for _, formatted := range w.buffer {
		if _, err := w.file.Write(formatted); err != nil {
			return err
		}
	}

	w.buffer = w.buffer[:0] // バッファをクリア
	return w.file.Sync()
}

// Close はファイルを閉じ、バッファを出力する
func (w *BufferedFileWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.FileWriter.Close()
}
