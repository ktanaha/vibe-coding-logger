package main

import (
	"errors"
	"time"

	"vibe-coding-logger/pkg/logger"
)

func main() {
	// システム情報付きロガーを作成
	log := logger.Default()
	
	// システム情報とランタイム情報を有効化
	log.EnableSystemInfo(true)
	log.EnableRuntimeInfo(true)
	
	log.Info("システム情報ロギングのデモを開始")
	
	// 基本的なシステム情報の表示
	systemInfo := logger.GetSystemInfo()
	log.Info("システム情報を取得", 
		logger.String("os", systemInfo.OS),
		logger.String("architecture", systemInfo.Architecture),
		logger.String("go_version", systemInfo.GoVersion),
		logger.String("hostname", systemInfo.Hostname),
		logger.Int("num_cpu", systemInfo.NumCPU))

	// 環境情報の表示
	envInfo := logger.GetEnvironmentInfo()
	log.Info("環境情報を取得",
		logger.String("working_directory", envInfo.WorkingDirectory),
		logger.String("go_path", envInfo.GoPath),
		logger.String("git_branch", envInfo.GitBranch),
		logger.String("git_commit", envInfo.GitCommit))

	// コンパクトなシステム情報を使用したログ
	compactInfo := logger.GetCompactSystemInfo()
	log.Info("コンパクトシステム情報", logger.Any("compact_info", compactInfo))

	// バイブコーディングセッションでシステム情報を記録
	sessionID := "system_demo_session"
	problemDomain := "logging_library_development"
	
	vibeTracker := logger.NewVibeTracker(log, sessionID, problemDomain, "development")
	
	// 開発環境の詳細を記録
	vibeTracker.LogThinkingProcess(
		"ロギングライブラリにシステム情報機能を追加",
		[]string{
			"OSとハードウェア情報の自動収集",
			"Git情報の取得",
			"言語バージョンの記録",
			"パフォーマンスへの影響を最小化",
		})

	// パフォーマンステストのシミュレーション
	start := time.Now()
	
	// 大量のログ出力（システム情報付き）
	for i := 0; i < 100; i++ {
		log.Debug("パフォーマンステスト", 
			logger.Int("iteration", i),
			logger.String("test_type", "system_info_logging"))
	}
	
	duration := time.Since(start)
	
	// パフォーマンス結果を記録
	vibeTracker.LogTestResult(
		"SystemInfoLoggingPerformance",
		true,
		"100回のログ出力を完了",
		duration)

	// ランタイム統計を記録
	runtimeStats := logger.GetRuntimeStats()
	log.Info("ランタイム統計", logger.Any("runtime_stats", runtimeStats))

	// システム情報を無効化してパフォーマンス比較
	log.EnableSystemInfo(false)
	log.EnableRuntimeInfo(false)
	
	start2 := time.Now()
	
	// システム情報なしでの大量ログ出力
	for i := 0; i < 100; i++ {
		log.Debug("パフォーマンステスト（システム情報なし）", 
			logger.Int("iteration", i),
			logger.String("test_type", "basic_logging"))
	}
	
	duration2 := time.Since(start2)
	
	// 再度システム情報を有効化して結果を記録
	log.EnableSystemInfo(true)
	
	log.Info("パフォーマンス比較結果",
		logger.Duration("with_system_info", duration),
		logger.Duration("without_system_info", duration2),
		logger.Any("performance_ratio", float64(duration.Nanoseconds())/float64(duration2.Nanoseconds())))

	// エラーハンドリング時のシステム情報
	errorHandler := logger.NewVibeErrorHandler(log, sessionID, problemDomain)
	
	testError := errors.New("システム情報収集エラー")
	errorHandler.HandleCodingError(
		testError,
		"system_info.go",
		125,
		"systemInfo := collectSystemInfo()",
		"エラーハンドリングを改善し、フォールバック機能を追加")

	// セッション終了時のサマリー
	vibeTracker.LogSessionSummary(
		[]string{
			"システム情報自動収集機能の実装完了",
			"パフォーマンス影響の測定完了",
			"エラーハンドリングの改善",
		},
		[]string{
			"システム情報収集のオーバーヘッド",
			"異なるOS間での互換性",
		},
		[]string{
			"キャッシュ機能による性能向上",
			"設定可能な情報収集レベル",
			"デバッグ時の環境情報の重要性",
		},
		[]string{
			"Git情報取得の改善",
			"Docker環境での情報収集",
			"CI/CD環境での特別な情報収集",
		})

	log.Info("システム情報ロギングのデモを完了",
		logger.Bool("system_info_enabled", log.IsSystemInfoEnabled()),
		logger.Bool("runtime_info_enabled", log.IsRuntimeInfoEnabled()))
}