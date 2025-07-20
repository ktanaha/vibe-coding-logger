package logger

import (
	"vibe-coding-logger/internal"
	"vibe-coding-logger/internal/formatter"
	"vibe-coding-logger/internal/writer"
)

// convertToInternalEntry はpkg/loggerのEntryをinternal.Entryに変換します
func convertToInternalEntry(entry *Entry) *internal.Entry {
	return &internal.Entry{
		ID:          entry.ID,
		Timestamp:   entry.Timestamp,
		Level:       internal.LogLevel(entry.Level),
		Action:      internal.ActionType(entry.Action),
		Operation:   entry.Operation,
		Input:       entry.Input,
		Output:      entry.Output,
		Error:       convertToInternalErrorInfo(entry.Error),
		Duration:    entry.Duration,
		Context:     entry.Context,
		Tags:        entry.Tags,
		TraceID:     entry.TraceID,
		SpanID:      entry.SpanID,
		ParentID:    entry.ParentID,
		Metadata:    entry.Metadata,
		SystemInfo:  entry.SystemInfo,
		RuntimeInfo: entry.RuntimeInfo,
	}
}

// convertToInternalErrorInfo はpkg/loggerのErrorInfoをinternal.ErrorInfoに変換します
func convertToInternalErrorInfo(errorInfo *ErrorInfo) *internal.ErrorInfo {
	if errorInfo == nil {
		return nil
	}
	return &internal.ErrorInfo{
		Message:    errorInfo.Message,
		Type:       errorInfo.Type,
		Code:       errorInfo.Code,
		Stack:      errorInfo.Stack,
		Retryable:  errorInfo.Retryable,
		Resolution: errorInfo.Resolution,
		Context:    errorInfo.Context,
	}
}

// NewConsoleWriter は新しいコンソールライターを作成します
func NewConsoleWriter() Writer {
	return &consoleWriterImpl{
		internalWriter: writer.NewConsoleWriter(),
	}
}

// NewTextFormatter は新しいテキストフォーマッターを作成します
func NewTextFormatter() Formatter {
	return &formatterAdapter{
		internalFormatter: formatter.NewTextFormatter(),
	}
}

// NewFileWriter は新しいファイルライターを作成します
func NewFileWriter(filename string) (Writer, error) {
	internalWriter, err := writer.NewFileWriter(filename)
	if err != nil {
		return nil, err
	}
	return &fileWriterImpl{
		internalWriter: internalWriter,
	}, nil
}

// consoleWriterImpl はinternal/writerを使用したWriter実装
type consoleWriterImpl struct {
	internalWriter internal.Writer
}

func (cw *consoleWriterImpl) Write(entry *Entry) error {
	internalEntry := convertToInternalEntry(entry)
	return cw.internalWriter.Write(internalEntry)
}

func (cw *consoleWriterImpl) Close() error {
	return cw.internalWriter.Close()
}

// formatterAdapter はinternal.FormatterをLogger.Formatterにアダプトします
type formatterAdapter struct {
	internalFormatter internal.Formatter
}

func (fa *formatterAdapter) Format(entry *Entry) ([]byte, error) {
	internalEntry := convertToInternalEntry(entry)
	return fa.internalFormatter.Format(internalEntry)
}

// fileWriterImpl はinternal/writerを使用したFileWriter実装
type fileWriterImpl struct {
	internalWriter *writer.FileWriter
}

func (fw *fileWriterImpl) Write(entry *Entry) error {
	internalEntry := convertToInternalEntry(entry)
	return fw.internalWriter.Write(internalEntry)
}

func (fw *fileWriterImpl) Close() error {
	return fw.internalWriter.Close()
}

func (fw *fileWriterImpl) SetFormatter(formatter Formatter) {
	if fa, ok := formatter.(*formatterAdapter); ok {
		fw.internalWriter.SetFormatter(fa.internalFormatter)
	}
}

// GetInternalFormatter は内部フォーマッターを取得します
func (fa *formatterAdapter) GetInternalFormatter() internal.Formatter {
	return fa.internalFormatter
}
