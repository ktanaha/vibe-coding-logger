package main

import (
	"errors"
	"time"

	"vibe-coding-logger/pkg/logger"
)

func main() {
	// デフォルトのロガーを作成
	log := logger.Default()

	// 基本的なログ出力
	log.Info("アプリケーションを開始します")
	log.Debug("デバッグメッセージ", logger.String("component", "main"))

	// フィールド付きログ
	log.Info("ユーザーログイン", 
		logger.String("user_id", "user123"),
		logger.String("ip_address", "192.168.1.100"),
		logger.Int("login_attempt", 1))

	// コンテキスト付きロガー
	userLogger := log.WithField("user_id", "user123").WithTag("authentication")
	userLogger.Info("認証成功")

	// エラーログ
	err := errors.New("データベース接続エラー")
	log.LogError(err, map[string]interface{}{
		"database": "postgres",
		"host":     "localhost",
		"port":     5432,
	}, true)

	// 操作追跡の例
	tracker := log.StartOperation("user_registration", map[string]interface{}{
		"email":    "user@example.com",
		"username": "newuser",
	})

	// 何らかの処理をシミュレート
	time.Sleep(100 * time.Millisecond)

	// 操作完了
	log.CompleteOperation(tracker, map[string]interface{}{
		"user_id": "user123",
		"status":  "active",
	})

	log.Info("アプリケーションを終了します")
}