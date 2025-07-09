// Package logger は操作のトラッキングとログ記録を行うためのパッケージです。
// バイブコーディングセッションの記録、一般的な操作の追跡、バッチ操作の管理を提供します。
package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Complete は操作の完了を記録します。
// 操作の実行結果をLoggerに渡し、完了状態として記録します。
func (t *OperationTracker) Complete(output map[string]interface{}) {
	t.Logger.CompleteOperation(t, output)
}

// Error は操作のエラーを記録します。
// 発生したエラーと、そのエラーに対する対処方法をLoggerに渡します。
func (t *OperationTracker) Error(err error, resolution string) {
	t.Logger.ErrorOperation(t, err, resolution)
}

// GetDuration は操作の経過時間を取得します。
// 操作開始時刻からの経過時間を計算して返します。
func (t *OperationTracker) GetDuration() time.Duration {
	return time.Since(t.StartTime)
}

// AddContext は操作にコンテキスト情報を追加します。
// 操作に関連する追加情報をキーと値のペアで保存します。
// Contextマップがnilの場合は初期化します。
func (t *OperationTracker) AddContext(key string, value interface{}) {
	if t.Context == nil {
		t.Context = make(map[string]interface{})
	}
	t.Context[key] = value
}

// CreateSubOperation は子操作を作成します。
// 親操作から派生した子操作を作成し、親のコンテキストを継承します。
// 子操作の開始ログを記録し、親子関係を管理します。
func (t *OperationTracker) CreateSubOperation(operation string, input map[string]interface{}) *OperationTracker {
	subTracker := &OperationTracker{
		ID:        uuid.New().String(),
		Operation: operation,
		StartTime: time.Now(),
		Input:     input,
		Context:   make(map[string]interface{}),
		Logger:    t.Logger,
		parent:    t,
	}

	// 親のコンテキストを継承
	for k, v := range t.Context {
		subTracker.Context[k] = v
	}

	t.Logger.Info(operation,
		String("action", string(ActionStart)),
		String("operation_id", subTracker.ID),
		String("parent_id", t.ID),
		Any("input", input))

	return subTracker
}

// VibeTracker はバイブコーディングセッション専用のトラッカーです。
// 一般的なOperationTrackerを組み込み、セッション管理と
// 問題領域、プログラミングステップに特化した機能を提供します。
type VibeTracker struct {
	*OperationTracker                        // 基本の操作トラッカー機能を組み込み
	sessionID         string                 // セッションの一意識別子
	problemDomain     string                 // 問題領域（例：「Web開発」「データ分析」）
	programmingStep   string                 // プログラミングステップ（例：「設計」「実装」「テスト」）
	contextData       map[string]interface{} // バイブコーディング固有のコンテキストデータ
}

// NewVibeTracker は新しいバイブトラッカーを作成します。
// セッションID、問題領域、プログラミングステップを設定し、
// 環境情報のスナップショットを自動的に記録します。
func NewVibeTracker(logger Logger, sessionID, problemDomain, programmingStep string) *VibeTracker {
	baseTracker := &OperationTracker{
		ID:        uuid.New().String(),
		Operation: fmt.Sprintf("vibe_coding_%s", programmingStep),
		StartTime: time.Now(),
		Input:     make(map[string]interface{}),
		Context:   make(map[string]interface{}),
		Logger:    logger,
	}

	vt := &VibeTracker{
		OperationTracker: baseTracker,
		sessionID:        sessionID,
		problemDomain:    problemDomain,
		programmingStep:  programmingStep,
		contextData:      make(map[string]interface{}),
	}

	// バイブコーディング固有のコンテキストを設定
	vt.contextData["session_id"] = sessionID
	vt.contextData["problem_domain"] = problemDomain
	vt.contextData["programming_step"] = programmingStep

	// システム環境情報を記録
	vt.LogEnvironmentSnapshot()

	return vt
}

// LogEnvironmentSnapshot は現在の環境のスナップショットを記録します。
// システム情報と環境情報を収集し、セッション開始時の状態を記録します。
func (vt *VibeTracker) LogEnvironmentSnapshot() {
	systemInfo := GetCompactSystemInfo()
	envInfo := GetEnvironmentInfo()

	vt.Logger.Info("environment_snapshot",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		Any("system_info", systemInfo),
		Any("environment_info", map[string]interface{}{
			"working_directory": envInfo.WorkingDirectory,
			"go_path":           envInfo.GoPath,
			"go_root":           envInfo.GoRoot,
			"go_mod":            envInfo.GoMod,
			"editor":            envInfo.Editor,
			"git_branch":        envInfo.GitBranch,
			"git_commit":        envInfo.GitCommit,
			"git_repository":    envInfo.GitRepository,
			"node_version":      envInfo.NodeVersion,
			"python_version":    envInfo.PythonVersion,
			"docker_version":    envInfo.DockerVersion,
		}))
}

// LogThinkingProcess は思考プロセスを記録します。
// 開発者の考え方や考慮事項を記録し、バイブコーディングの思考過程を追跡します。
func (vt *VibeTracker) LogThinkingProcess(thoughts string, considerations []string) {
	vt.Logger.Info("thinking_process",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("thoughts", thoughts),
		Any("considerations", considerations))
}

// LogDecision は決定を記録します。
// 設計上の決定、その理由、検討した代替案を記録します。
func (vt *VibeTracker) LogDecision(decision string, reasoning string, alternatives []string) {
	vt.Logger.Info("decision_made",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("decision", decision),
		String("reasoning", reasoning),
		Any("alternatives", alternatives))
}

// LogCodeChange はコード変更を記録します。
// ファイルの変更内容、変更理由、変更前後のコードを記録します。
func (vt *VibeTracker) LogCodeChange(filename string, changeType string, beforeCode string, afterCode string, reason string) {
	vt.Logger.Info("code_change",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("filename", filename),
		String("change_type", changeType),
		String("before_code", beforeCode),
		String("after_code", afterCode),
		String("reason", reason))
}

// LogTestResult はテスト結果を記録します。
// テストの実行結果、出力、実行時間を記録し、失敗時はERRORレベルでログを出力します。
func (vt *VibeTracker) LogTestResult(testName string, passed bool, output string, duration time.Duration) {
	logLevel := INFO
	if !passed {
		logLevel = ERROR
	}

	if logLevel == ERROR {
		vt.Logger.Error("test_result",
			String("session_id", vt.sessionID),
			String("problem_domain", vt.problemDomain),
			String("programming_step", vt.programmingStep),
			String("test_name", testName),
			String("result", map[bool]string{true: "PASSED", false: "FAILED"}[passed]),
			String("output", output),
			Duration("duration", duration))
	}
}

// LogRefactoring はリファクタリングを記録します。
// リファクタリングの種類、対象、理由、変更内容を記録します。
func (vt *VibeTracker) LogRefactoring(refactorType string, target string, reason string, before string, after string) {
	vt.Logger.Info("refactoring",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("refactor_type", refactorType),
		String("target", target),
		String("reason", reason),
		String("before", before),
		String("after", after))
}

// LogDebugSession はデバッグセッションを記録します。
// 問題、仮説、調査内容、解決方法を記録し、デバッグプロセスを追跡します。
func (vt *VibeTracker) LogDebugSession(issue string, hypothesis string, investigation string, resolution string) {
	vt.Logger.Info("debug_session",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("issue", issue),
		String("hypothesis", hypothesis),
		String("investigation", investigation),
		String("resolution", resolution))
}

// LogLearning は学習内容を記録します。
// 学習した概念、理解内容、適用方法、メモを記録し、知識の蓄積を追跡します。
func (vt *VibeTracker) LogLearning(concept string, understanding string, application string, notes string) {
	vt.Logger.Info("learning",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("concept", concept),
		String("understanding", understanding),
		String("application", application),
		String("notes", notes))
}

// LogBlocker はブロッカーを記録します。
// 開発を阻害する要因、影響、回避策、解決方法をWARNレベルで記録します。
func (vt *VibeTracker) LogBlocker(blocker string, impact string, workaround string, resolution string) {
	vt.Logger.Warn("blocker",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("blocker", blocker),
		String("impact", impact),
		String("workaround", workaround),
		String("resolution", resolution))
}

// LogBreakthrough はブレイクスルーを記録します。
// 重要な発見、コンテキスト、影響、学んだ教訓を記録します。
func (vt *VibeTracker) LogBreakthrough(breakthrough string, context string, impact string, lessons []string) {
	vt.Logger.Info("breakthrough",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("breakthrough", breakthrough),
		String("context", context),
		String("impact", impact),
		Any("lessons", lessons))
}

// LogSessionSummary はセッションサマリーを記録します。
// セッションの成果、課題、洞察、次のステップ、セッション時間を記録します。
func (vt *VibeTracker) LogSessionSummary(accomplishments []string, challenges []string, insights []string, nextSteps []string) {
	vt.Logger.Info("session_summary",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		Any("accomplishments", accomplishments),
		Any("challenges", challenges),
		Any("insights", insights),
		Any("next_steps", nextSteps),
		Duration("session_duration", vt.GetDuration()))
}

// BatchOperationTracker は複数の操作を一括で追跡します。
// 関連する複数の操作をグループ化し、バッチ全体の結果を管理します。
type BatchOperationTracker struct {
	ID         string                 // バッチの一意識別子
	BatchName  string                 // バッチの名前
	Operations []*OperationTracker    // バッチに含まれる操作のスライス
	StartTime  time.Time              // バッチの開始時刻
	Logger     Logger                 // ログ出力用のLoggerインスタンス
	Context    map[string]interface{} // バッチ全体のコンテキスト情報
}

// NewBatchOperationTracker は新しいバッチ操作トラッカーを作成します。
// バッチ名とLoggerを指定し、初期化されたバッチトラッカーを返します。
func NewBatchOperationTracker(logger Logger, batchName string) *BatchOperationTracker {
	return &BatchOperationTracker{
		ID:         uuid.New().String(),
		BatchName:  batchName,
		Operations: make([]*OperationTracker, 0),
		StartTime:  time.Now(),
		Logger:     logger,
		Context:    make(map[string]interface{}),
	}
}

// AddOperation はバッチに操作を追加します。
// 新しい操作を作成し、バッチのコンテキストを継承し、操作リストに追加します。
func (bt *BatchOperationTracker) AddOperation(operation string, input map[string]interface{}) *OperationTracker {
	tracker := &OperationTracker{
		ID:        uuid.New().String(),
		Operation: operation,
		StartTime: time.Now(),
		Input:     input,
		Context:   make(map[string]interface{}),
		Logger:    bt.Logger,
	}

	// バッチのコンテキストを継承
	for k, v := range bt.Context {
		tracker.Context[k] = v
	}

	bt.Operations = append(bt.Operations, tracker)

	bt.Logger.Info(operation,
		String("action", string(ActionStart)),
		String("operation_id", tracker.ID),
		String("batch_id", bt.ID),
		String("batch_name", bt.BatchName),
		Any("input", input))

	return tracker
}

// Complete はバッチ操作の完了を記録します。
// バッチ全体の結果、統計情報、サマリーを記録します。
func (bt *BatchOperationTracker) Complete(summary map[string]interface{}) {
	duration := time.Since(bt.StartTime)

	stats := map[string]interface{}{
		"total_operations":     len(bt.Operations),
		"completed_operations": bt.countCompletedOperations(),
		"failed_operations":    bt.countFailedOperations(),
		"total_duration":       duration,
	}

	bt.Logger.Info(bt.BatchName,
		String("action", string(ActionComplete)),
		String("batch_id", bt.ID),
		Any("summary", summary),
		Any("stats", stats))
}

// countCompletedOperations は完了した操作の数を数えます。
// 現在の実装ではすべての操作を完了とみなします。
func (bt *BatchOperationTracker) countCompletedOperations() int {
	// TODO: 実際の実装では、各操作の状態を追跡する必要があります
	return len(bt.Operations)
}

// countFailedOperations は失敗した操作の数を数えます。
// 現在の実装では失敗した操作はないとみなします。
func (bt *BatchOperationTracker) countFailedOperations() int {
	// TODO: 実際の実装では、各操作の状態を追跡する必要があります
	return 0
}

// LogOperationMetrics は操作メトリクスを記録します。
// 操作の結果、実行時間、メタデータを記録し、パフォーマンス分析を支援します。
func LogOperationMetrics(ctx context.Context, logger Logger, operation string, duration time.Duration, success bool, metadata map[string]interface{}) {
	contextLogger := logger.WithContext(ctx)
	if success {
		contextLogger.Info("operation_metrics",
			String("operation", operation),
			String("result", "SUCCESS"),
			Duration("duration", duration),
			Any("metadata", metadata))
	} else {
		contextLogger.Error("operation_metrics",
			String("operation", operation),
			String("result", "FAILURE"),
			Duration("duration", duration),
			Any("metadata", metadata))
	}
}
