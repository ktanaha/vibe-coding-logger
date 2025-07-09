package logger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SystemInfo はシステム情報を表す
type SystemInfo struct {
	// 基本システム情報
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Hostname     string `json:"hostname"`
	
	// Go runtime情報
	GoVersion    string `json:"go_version"`
	GoOS         string `json:"go_os"`
	GoArch       string `json:"go_arch"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	
	// メモリ情報
	MemStats     *runtime.MemStats `json:"mem_stats,omitempty"`
	
	// プロセス情報
	PID          int    `json:"pid"`
	
	// 環境変数（選択的）
	Environment  map[string]string `json:"environment,omitempty"`
	
	// タイムスタンプ
	CollectedAt  time.Time `json:"collected_at"`
}

// EnvironmentInfo は環境情報を表す
type EnvironmentInfo struct {
	// 開発環境情報
	WorkingDirectory string `json:"working_directory"`
	GoPath          string `json:"go_path"`
	GoRoot          string `json:"go_root"`
	GoMod           string `json:"go_mod,omitempty"`
	
	// エディタ・IDE情報
	Editor          string `json:"editor,omitempty"`
	
	// Git情報
	GitBranch       string `json:"git_branch,omitempty"`
	GitCommit       string `json:"git_commit,omitempty"`
	GitRepository   string `json:"git_repository,omitempty"`
	
	// バージョン管理
	NodeVersion     string `json:"node_version,omitempty"`
	PythonVersion   string `json:"python_version,omitempty"`
	DockerVersion   string `json:"docker_version,omitempty"`
	
	// その他のツール
	DatabaseVersion map[string]string `json:"database_version,omitempty"`
	
	CollectedAt     time.Time `json:"collected_at"`
}

// SystemInfoCollector はシステム情報を収集する
type SystemInfoCollector struct {
	mu    sync.RWMutex
	cache *SystemInfo
	envCache *EnvironmentInfo
	cacheExpiry time.Duration
	lastCollected time.Time
}

// NewSystemInfoCollector は新しいシステム情報コレクターを作成する
func NewSystemInfoCollector() *SystemInfoCollector {
	return &SystemInfoCollector{
		cacheExpiry: 5 * time.Minute, // 5分間キャッシュ
	}
}

// GetSystemInfo は現在のシステム情報を取得する
func (sic *SystemInfoCollector) GetSystemInfo() *SystemInfo {
	sic.mu.RLock()
	if sic.cache != nil && time.Since(sic.lastCollected) < sic.cacheExpiry {
		defer sic.mu.RUnlock()
		return sic.cache
	}
	sic.mu.RUnlock()

	sic.mu.Lock()
	defer sic.mu.Unlock()

	// ダブルチェック
	if sic.cache != nil && time.Since(sic.lastCollected) < sic.cacheExpiry {
		return sic.cache
	}

	sic.cache = sic.collectSystemInfo()
	sic.lastCollected = time.Now()
	return sic.cache
}

// GetEnvironmentInfo は環境情報を取得する
func (sic *SystemInfoCollector) GetEnvironmentInfo() *EnvironmentInfo {
	sic.mu.RLock()
	if sic.envCache != nil && time.Since(sic.lastCollected) < sic.cacheExpiry {
		defer sic.mu.RUnlock()
		return sic.envCache
	}
	sic.mu.RUnlock()

	sic.mu.Lock()
	defer sic.mu.Unlock()

	sic.envCache = sic.collectEnvironmentInfo()
	return sic.envCache
}

// collectSystemInfo はシステム情報を収集する
func (sic *SystemInfoCollector) collectSystemInfo() *SystemInfo {
	hostname, _ := os.Hostname()
	
	// メモリ統計を取得
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 重要な環境変数を収集
	importantEnvVars := []string{
		"PATH", "HOME", "USER", "SHELL", "TERM",
		"GOPATH", "GOROOT", "GOOS", "GOARCH",
		"CI", "DOCKER", "KUBERNETES_SERVICE_HOST",
	}
	
	environment := make(map[string]string)
	for _, key := range importantEnvVars {
		if value := os.Getenv(key); value != "" {
			environment[key] = value
		}
	}

	return &SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Hostname:     hostname,
		GoVersion:    runtime.Version(),
		GoOS:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemStats:     &memStats,
		PID:          os.Getpid(),
		Environment:  environment,
		CollectedAt:  time.Now(),
	}
}

// collectEnvironmentInfo は環境情報を収集する
func (sic *SystemInfoCollector) collectEnvironmentInfo() *EnvironmentInfo {
	wd, _ := os.Getwd()
	
	envInfo := &EnvironmentInfo{
		WorkingDirectory: wd,
		GoPath:          os.Getenv("GOPATH"),
		GoRoot:          os.Getenv("GOROOT"),
		CollectedAt:     time.Now(),
	}

	// go.modファイルの存在確認
	if _, err := os.Stat("go.mod"); err == nil {
		envInfo.GoMod = "present"
	}

	// エディタ情報（環境変数から推測）
	if editor := os.Getenv("EDITOR"); editor != "" {
		envInfo.Editor = editor
	} else if visual := os.Getenv("VISUAL"); visual != "" {
		envInfo.Editor = visual
	}

	// Git情報を収集（エラーは無視）
	envInfo.GitBranch = sic.getGitBranch()
	envInfo.GitCommit = sic.getGitCommit()
	envInfo.GitRepository = sic.getGitRepository()

	// その他の言語バージョン情報
	envInfo.NodeVersion = sic.getCommandVersion("node", "--version")
	envInfo.PythonVersion = sic.getCommandVersion("python", "--version")
	envInfo.DockerVersion = sic.getDockerVersion()

	// データベースバージョン情報
	envInfo.DatabaseVersion = sic.getDatabaseVersions()

	return envInfo
}

// getGitBranch は現在のGitブランチを取得する
func (sic *SystemInfoCollector) getGitBranch() string {
	// .git/HEADファイルから読み取り
	if data, err := os.ReadFile(".git/HEAD"); err == nil {
		content := strings.TrimSpace(string(data))
		if strings.HasPrefix(content, "ref: refs/heads/") {
			return strings.TrimPrefix(content, "ref: refs/heads/")
		}
	}
	return ""
}

// getGitCommit は現在のGitコミットハッシュを取得する
func (sic *SystemInfoCollector) getGitCommit() string {
	// .git/HEADファイルから読み取り
	if data, err := os.ReadFile(".git/HEAD"); err == nil {
		content := strings.TrimSpace(string(data))
		if strings.HasPrefix(content, "ref: ") {
			// ブランチの場合、そのブランチのコミットを読み取り
			refPath := strings.TrimPrefix(content, "ref: ")
			if commitData, err := os.ReadFile(".git/" + refPath); err == nil {
				return strings.TrimSpace(string(commitData))[:8] // 短縮ハッシュ
			}
		} else if len(content) >= 8 {
			// 直接コミットハッシュの場合
			return content[:8]
		}
	}
	return ""
}

// getGitRepository はGitリポジトリ情報を取得する
func (sic *SystemInfoCollector) getGitRepository() string {
	// .git/configファイルからremote origin URLを読み取り
	if data, err := os.ReadFile(".git/config"); err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		inOriginSection := false
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == `[remote "origin"]` {
				inOriginSection = true
				continue
			}
			if strings.HasPrefix(line, "[") && inOriginSection {
				break
			}
			if inOriginSection && strings.HasPrefix(line, "url = ") {
				return strings.TrimPrefix(line, "url = ")
			}
		}
	}
	return ""
}

// getCommandVersion は指定したコマンドのバージョンを取得する
func (sic *SystemInfoCollector) getCommandVersion(command, versionFlag string) string {
	// セキュリティ上の理由で、実際の実装では外部コマンド実行は避ける
	// ここでは環境変数やファイルベースの検出に留める
	
	switch command {
	case "node":
		if version := os.Getenv("NODE_VERSION"); version != "" {
			return version
		}
	case "python":
		if version := os.Getenv("PYTHON_VERSION"); version != "" {
			return version
		}
	}
	
	return ""
}

// getDockerVersion はDockerのバージョン情報を取得する
func (sic *SystemInfoCollector) getDockerVersion() string {
	// Docker環境の検出
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "running-in-container"
	}
	
	if version := os.Getenv("DOCKER_VERSION"); version != "" {
		return version
	}
	
	return ""
}

// getDatabaseVersions はデータベースのバージョン情報を取得する
func (sic *SystemInfoCollector) getDatabaseVersions() map[string]string {
	versions := make(map[string]string)
	
	// 環境変数からデータベース情報を収集
	dbEnvVars := map[string]string{
		"POSTGRES_VERSION": "postgresql",
		"MYSQL_VERSION":    "mysql",
		"REDIS_VERSION":    "redis",
		"MONGODB_VERSION":  "mongodb",
	}
	
	for envVar, dbName := range dbEnvVars {
		if version := os.Getenv(envVar); version != "" {
			versions[dbName] = version
		}
	}
	
	return versions
}

// GetCompactSystemInfo はコンパクトなシステム情報を取得する
func (sic *SystemInfoCollector) GetCompactSystemInfo() map[string]interface{} {
	sysInfo := sic.GetSystemInfo()
	envInfo := sic.GetEnvironmentInfo()
	
	compact := map[string]interface{}{
		"os":           sysInfo.OS,
		"arch":         sysInfo.Architecture,
		"go_version":   sysInfo.GoVersion,
		"hostname":     sysInfo.Hostname,
		"pid":          sysInfo.PID,
		"num_cpu":      sysInfo.NumCPU,
		"working_dir":  envInfo.WorkingDirectory,
	}
	
	// Git情報があれば追加
	if envInfo.GitBranch != "" {
		compact["git_branch"] = envInfo.GitBranch
	}
	if envInfo.GitCommit != "" {
		compact["git_commit"] = envInfo.GitCommit
	}
	
	// その他の言語バージョン
	if envInfo.NodeVersion != "" {
		compact["node_version"] = envInfo.NodeVersion
	}
	if envInfo.PythonVersion != "" {
		compact["python_version"] = envInfo.PythonVersion
	}
	
	return compact
}

// GetRuntimeStats はランタイム統計を取得する
func (sic *SystemInfoCollector) GetRuntimeStats() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return map[string]interface{}{
		"goroutines":     runtime.NumGoroutine(),
		"heap_alloc":     memStats.HeapAlloc,
		"heap_sys":       memStats.HeapSys,
		"heap_objects":   memStats.HeapObjects,
		"stack_inuse":    memStats.StackInuse,
		"gc_runs":        memStats.NumGC,
		"next_gc":        memStats.NextGC,
		"last_gc":        time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339),
	}
}

// defaultSystemInfoCollector はデフォルトのシステム情報コレクター
var defaultSystemInfoCollector = NewSystemInfoCollector()

// GetSystemInfo はデフォルトコレクターからシステム情報を取得する
func GetSystemInfo() *SystemInfo {
	return defaultSystemInfoCollector.GetSystemInfo()
}

// GetEnvironmentInfo はデフォルトコレクターから環境情報を取得する
func GetEnvironmentInfo() *EnvironmentInfo {
	return defaultSystemInfoCollector.GetEnvironmentInfo()
}

// GetCompactSystemInfo はデフォルトコレクターからコンパクトなシステム情報を取得する
func GetCompactSystemInfo() map[string]interface{} {
	return defaultSystemInfoCollector.GetCompactSystemInfo()
}

// GetRuntimeStats はデフォルトコレクターからランタイム統計を取得する
func GetRuntimeStats() map[string]interface{} {
	return defaultSystemInfoCollector.GetRuntimeStats()
}