package main

import (
	"errors"
	"time"

	"vibe-coding-logger/pkg/logger"
)

func main() {
	// バイブコーディング専用ロガーを作成
	log := logger.Default()
	
	// セッション情報を設定
	sessionID := "vibe_session_2025_01_07"
	problemDomain := "rest_api_development"
	
	// バイブトラッカーを作成
	vibeTracker := logger.NewVibeTracker(log, sessionID, problemDomain, "planning")

	// 思考プロセスを記録
	vibeTracker.LogThinkingProcess(
		"RESTful APIの設計について考える",
		[]string{
			"エンドポイントの設計",
			"認証方式の選択", 
			"エラーハンドリングの方針",
			"バリデーションの実装",
		})

	// 決定を記録
	vibeTracker.LogDecision(
		"JWT認証を採用",
		"スケーラビリティとステートレス性を重視",
		[]string{"セッション認証", "OAuth2"})

	// プログラミングステップを変更
	codingTracker := logger.NewVibeTracker(log, sessionID, problemDomain, "coding")

	// コード変更を記録
	codingTracker.LogCodeChange(
		"handlers/auth.go",
		"新規作成",
		"",
		`func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// JWT認証の実装
	token, err := generateJWT(userID)
	if err != nil {
		http.Error(w, "認証エラー", http.StatusInternalServerError)
		return
	}
	// ...
}`,
		"JWT認証のログインハンドラーを実装")

	// テスト結果を記録
	codingTracker.LogTestResult(
		"TestLoginHandler",
		true,
		"全てのテストケースが成功",
		50*time.Millisecond)

	// エラーハンドリングの例
	errorHandler := logger.NewVibeErrorHandler(log, sessionID, problemDomain)
	
	testError := errors.New("データベース接続タイムアウト")
	errorHandler.HandleCodingError(
		testError,
		"models/user.go",
		45,
		"db.Query(sql, params...)",
		"接続プールの設定を見直し、タイムアウト値を調整")

	// リファクタリングを記録
	refactorTracker := logger.NewVibeTracker(log, sessionID, problemDomain, "refactoring")
	
	refactorTracker.LogRefactoring(
		"Extract Method",
		"LoginHandler内の認証ロジック",
		"メソッドが長すぎて理解しにくい",
		"長いLoginHandlerメソッド",
		"authenticateUser()とgenerateToken()に分離")

	// デバッグセッションを記録
	refactorTracker.LogDebugSession(
		"JWTトークンの有効期限が正しく設定されない",
		"time.Duration の設定に問題がある可能性",
		"ログ出力とブレークポイントで値を確認",
		"time.Hour * 24 を time.Hour * 24 * 7 に修正")

	// 学習内容を記録
	refactorTracker.LogLearning(
		"JWT仕様の詳細",
		"クレームの種類とセキュリティベストプラクティス",
		"expires_inクレームとsub（subject）クレームの適切な使用",
		"RFC 7519をさらに詳しく読む必要がある")

	// ブレイクスルーを記録
	refactorTracker.LogBreakthrough(
		"認証フローの完全な理解",
		"JWTとRefresh Tokenの組み合わせパターンの実装中",
		"セキュリティとユーザビリティの両立が可能",
		[]string{
			"アクセストークンの短期間設定",
			"リフレッシュトークンの安全な保存",
			"自動更新メカニズムの実装",
		})

	// セッションサマリーを記録
	refactorTracker.LogSessionSummary(
		[]string{
			"JWT認証システムの基本実装完了",
			"ユニットテストの作成とパス",
			"エラーハンドリングの改善",
			"コードリファクタリングによる保守性向上",
		},
		[]string{
			"データベース接続の安定性",
			"JWT仕様の複雑さ",
			"セキュリティ要件とパフォーマンスのバランス",
		},
		[]string{
			"認証フローの設計パターンの理解",
			"TDDによるコード品質の向上",
			"段階的リファクタリングの効果",
		},
		[]string{
			"リフレッシュトークン機能の実装",
			"認証ミドルウェアの作成",
			"API Rate Limitingの実装",
			"セキュリティテストの追加",
		})

	log.Info("バイブコーディングセッション完了", 
		logger.String("session_id", sessionID),
		logger.String("problem_domain", problemDomain),
		logger.Duration("total_duration", 2*time.Hour))
}