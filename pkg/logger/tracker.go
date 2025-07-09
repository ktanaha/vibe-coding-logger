package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Complete は操作の完了を記録する
func (t *OperationTracker) Complete(output map[string]interface{}) {
	t.Logger.CompleteOperation(t, output)
}

// Error は操作のエラーを記録する
func (t *OperationTracker) Error(err error, resolution string) {
	t.Logger.ErrorOperation(t, err, resolution)
}

// GetDuration は操作の経過時間を取得する
func (t *OperationTracker) GetDuration() time.Duration {
	return time.Since(t.StartTime)
}

// AddContext は操作にコンテキストを追加する
func (t *OperationTracker) AddContext(key string, value interface{}) {
	if t.Context == nil {
		t.Context = make(map[string]interface{})
	}
	t.Context[key] = value
}

// CreateSubOperation は子操作を作成する
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

	t.Logger.log(INFO, operation,
		String("action", string(ActionStart)),
		String("operation_id", subTracker.ID),
		String("parent_id", t.ID),
		Any("input", input))

	return subTracker
}

// VibeTracker はバイブコーディング専用のトラッカー
type VibeTracker struct {
	*OperationTracker
	sessionID       string
	problemDomain   string
	programmingStep string
	contextData     map[string]interface{}
}

// NewVibeTracker は新しいバイブトラッカーを作成する
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

// LogEnvironmentSnapshot は現在の環境のスナップショットを記録する
func (vt *VibeTracker) LogEnvironmentSnapshot() {
	systemInfo := GetCompactSystemInfo()
	envInfo := GetEnvironmentInfo()
	
	vt.Logger.Info( "environment_snapshot",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		Any("system_info", systemInfo),
		Any("environment_info", map[string]interface{}{
			"working_directory": envInfo.WorkingDirectory,
			"go_path":          envInfo.GoPath,
			"go_root":          envInfo.GoRoot,
			"go_mod":           envInfo.GoMod,
			"editor":           envInfo.Editor,
			"git_branch":       envInfo.GitBranch,
			"git_commit":       envInfo.GitCommit,
			"git_repository":   envInfo.GitRepository,
			"node_version":     envInfo.NodeVersion,
			"python_version":   envInfo.PythonVersion,
			"docker_version":   envInfo.DockerVersion,
		}))
}

// LogThinkingProcess は思考プロセスを記録する
func (vt *VibeTracker) LogThinkingProcess(thoughts string, considerations []string) {
	vt.Logger.Info("thinking_process",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("thoughts", thoughts),
		Any("considerations", considerations))
}

// LogDecision は決定を記録する
func (vt *VibeTracker) LogDecision(decision string, reasoning string, alternatives []string) {
	vt.Logger.Info( "decision_made",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("decision", decision),
		String("reasoning", reasoning),
		Any("alternatives", alternatives))
}

// LogCodeChange はコード変更を記録する
func (vt *VibeTracker) LogCodeChange(filename string, changeType string, beforeCode string, afterCode string, reason string) {
	vt.Logger.Info( "code_change",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("filename", filename),
		String("change_type", changeType),
		String("before_code", beforeCode),
		String("after_code", afterCode),
		String("reason", reason))
}

// LogTestResult はテスト結果を記録する
func (vt *VibeTracker) LogTestResult(testName string, passed bool, output string, duration time.Duration) {
	logLevel := INFO
	if !passed {
		logLevel = ERROR
	}

	if logLevel == ERROR {
		vt.Logger.Error( "test_result",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("test_name", testName),
		String("result", func() string {
			if passed {
				return "PASSED"
			}
			return "FAILED"
		}()),
		String("output", output),
		Duration("duration", duration))
}

// LogRefactoring はリファクタリングを記録する
func (vt *VibeTracker) LogRefactoring(refactorType string, target string, reason string, before string, after string) {
	vt.Logger.Info( "refactoring",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("refactor_type", refactorType),
		String("target", target),
		String("reason", reason),
		String("before", before),
		String("after", after))
}

// LogDebugSession はデバッグセッションを記録する
func (vt *VibeTracker) LogDebugSession(issue string, hypothesis string, investigation string, resolution string) {
	vt.Logger.Info( "debug_session",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("issue", issue),
		String("hypothesis", hypothesis),
		String("investigation", investigation),
		String("resolution", resolution))
}

// LogLearning は学習内容を記録する
func (vt *VibeTracker) LogLearning(concept string, understanding string, application string, notes string) {
	vt.Logger.Info( "learning",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("concept", concept),
		String("understanding", understanding),
		String("application", application),
		String("notes", notes))
}

// LogBlocker はブロッカーを記録する
func (vt *VibeTracker) LogBlocker(blocker string, impact string, workaround string, resolution string) {
	vt.Logger.log(WARN, "blocker",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("blocker", blocker),
		String("impact", impact),
		String("workaround", workaround),
		String("resolution", resolution))
}

// LogBreakthrough はブレイクスルーを記録する
func (vt *VibeTracker) LogBreakthrough(breakthrough string, context string, impact string, lessons []string) {
	vt.Logger.Info( "breakthrough",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		String("breakthrough", breakthrough),
		String("context", context),
		String("impact", impact),
		Any("lessons", lessons))
}

// LogSessionSummary はセッションサマリーを記録する
func (vt *VibeTracker) LogSessionSummary(accomplishments []string, challenges []string, insights []string, nextSteps []string) {
	vt.Logger.Info( "session_summary",
		String("session_id", vt.sessionID),
		String("problem_domain", vt.problemDomain),
		String("programming_step", vt.programmingStep),
		Any("accomplishments", accomplishments),
		Any("challenges", challenges),
		Any("insights", insights),
		Any("next_steps", nextSteps),
		Duration("session_duration", vt.GetDuration()))
}

// BatchOperationTracker は複数の操作を一括で追跡する
type BatchOperationTracker struct {
	ID         string
	BatchName  string
	Operations []*OperationTracker
	StartTime  time.Time
	Logger     Logger
	Context    map[string]interface{}
}

// NewBatchOperationTracker は新しいバッチ操作トラッカーを作成する
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

// AddOperation はバッチに操作を追加する
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

	bt.Logger.log(INFO, operation,
		String("action", string(ActionStart)),
		String("operation_id", tracker.ID),
		String("batch_id", bt.ID),
		String("batch_name", bt.BatchName),
		Any("input", input))

	return tracker
}

// Complete はバッチ操作の完了を記録する
func (bt *BatchOperationTracker) Complete(summary map[string]interface{}) {
	duration := time.Since(bt.StartTime)
	
	stats := map[string]interface{}{
		"total_operations": len(bt.Operations),
		"completed_operations": bt.countCompletedOperations(),
		"failed_operations": bt.countFailedOperations(),
		"total_duration": duration,
	}

	bt.Logger.log(INFO, bt.BatchName,
		String("action", string(ActionComplete)),
		String("batch_id", bt.ID),
		Any("summary", summary),
		Any("stats", stats))
}

// countCompletedOperations は完了した操作の数を数える
func (bt *BatchOperationTracker) countCompletedOperations() int {
	// この実装では簡単のため、すべての操作を完了とみなす
	// 実際の実装では、各操作の状態を追跡する必要がある
	return len(bt.Operations)
}

// countFailedOperations は失敗した操作の数を数える
func (bt *BatchOperationTracker) countFailedOperations() int {
	// この実装では簡単のため、失敗した操作はないとみなす
	// 実際の実装では、各操作の状態を追跡する必要がある
	return 0
}

// LogOperationMetrics は操作メトリクスを記録する
func LogOperationMetrics(ctx context.Context, logger Logger, operation string, duration time.Duration, success bool, metadata map[string]interface{}) {
	logLevel := INFO
	if !success {
		logLevel = ERROR
	}

	logger.WithContext(ctx).log(logLevel, "operation_metrics",
		String("operation", operation),
		String("result", func() string {
			if success {
				return "SUCCESS"
			}
			return "FAILURE"
		}()),
		Duration("duration", duration),
		Any("metadata", metadata))
}