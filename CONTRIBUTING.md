# 🤝 コントリビューションガイド

Vibe Coding Loggerへのコントリビューションを歓迎します！このガイドに従って、プロジェクトに貢献してください。

## 📋 貢献の種類

- 🐛 **バグ報告**: 問題を発見した場合
- ✨ **機能提案**: 新しい機能のアイデア
- 📝 **ドキュメント改善**: READMEや説明の改善
- 💻 **コード貢献**: バグ修正や新機能の実装
- 🧪 **テスト**: テストケースの追加や改善
- 🌍 **翻訳**: 多言語対応

## 🚀 開発環境のセットアップ

### 必要な環境

- Go 1.24+ 
- Git
- エディタ（VS Code、GoLandなど）

### セットアップ手順

1. **リポジトリのフォーク**
   ```bash
   # GitHubでフォークしてからクローン
   git clone https://github.com/your-username/vibe-coding-logger.git
   cd vibe-coding-logger
   ```

2. **依存関係のインストール**
   ```bash
   go mod download
   go mod tidy
   ```

3. **動作確認**
   ```bash
   # 基本機能のテスト
   go run simple_demo.go
   
   # 使用例の実行
   go run examples/basic_usage.go
   ```

## 🔄 開発ワークフロー

### 1. イシューの確認

- 既存のIssueを確認
- 新しい問題や提案がある場合はIssueを作成
- 大きな変更の場合は事前にDiscussionで相談

### 2. ブランチの作成

```bash
# メインブランチから最新を取得
git checkout main
git pull upstream main

# フィーチャーブランチを作成
git checkout -b feature/your-feature-name
# または
git checkout -b bugfix/issue-number
```

### 3. 開発とテスト

```bash
# 開発
# コードを編集...

# フォーマット
go fmt ./...

# リント（推奨）
golangci-lint run

# テスト
go test ./...

# ビルド確認
go build ./...
```

### 4. コミット

```bash
# 変更をステージング
git add .

# コミット（詳細は後述）
git commit -m "feat: 新機能の説明"

# プッシュ
git push origin feature/your-feature-name
```

### 5. プルリクエスト

1. GitHubでプルリクエストを作成
2. 適切なタイトルと説明を記述
3. レビューを待つ
4. フィードバックに対応

## 📝 コミットメッセージ規約

### フォーマット

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Type

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント変更
- `style`: コードフォーマット
- `refactor`: リファクタリング
- `test`: テスト関連
- `chore`: その他（依存関係更新など）

### 例

```bash
feat(logger): システム情報自動収集機能を追加

- OS、アーキテクチャ、Go バージョンを自動収集
- キャッシュ機能でパフォーマンスを最適化
- 有効/無効の切り替え機能を追加

Closes #123
```

## 🧪 テストガイドライン

### テストの種類

1. **ユニットテスト**: 個別関数のテスト
2. **統合テスト**: コンポーネント間の連携テスト
3. **パフォーマンステスト**: 性能テスト

### テストの実行

```bash
# 全テスト実行
go test ./...

# 詳細出力
go test -v ./...

# カバレッジ
go test -cover ./...

# ベンチマーク
go test -bench=. ./...
```

### テストの書き方

```go
func TestLoggerBasicFunctionality(t *testing.T) {
    // テスト準備
    log := logger.New(logger.INFO)
    
    // テスト実行
    log.Info("テストメッセージ")
    
    // 検証
    // アサーション...
}
```

## 📚 コーディング規約

### Go言語規約

- [Effective Go](https://golang.org/doc/effective_go.html)に従う
- `go fmt`でフォーマット
- `golint`でリント
- 公開関数には適切なドキュメントコメント

### プロジェクト固有規約

1. **ファイル構成**
   - `pkg/`: 公開API
   - `internal/`: 内部実装
   - `examples/`: 使用例
   - `tests/`: テストコード

2. **命名規約**
   - インターフェース: `Logger`, `Writer`, `Formatter`
   - 実装: `vibeLogger`, `fileWriter`, `jsonFormatter`
   - バイブコーディング用: `Vibe`プレフィックス

3. **エラーハンドリング**
   - 明示的なエラー処理
   - エラーメッセージは日本語対応
   - リカバリー機能の実装

### ドキュメント

```go
// LogThinkingProcess は思考プロセスを記録する
// 
// thoughts: 思考内容の説明
// considerations: 検討した選択肢のリスト
//
// 使用例:
//   vibeTracker.LogThinkingProcess(
//       "認証方式の検討", 
//       []string{"JWT", "Session"})
func (vt *VibeTracker) LogThinkingProcess(thoughts string, considerations []string) {
    // 実装...
}
```

## 🐛 バグ報告

### バグレポートに含める情報

- **環境情報**
  - Go バージョン
  - OS（バージョン含む）
  - アーキテクチャ

- **再現手順**
  ```
  1. ロガーを作成
  2. システム情報を有効化
  3. ログを出力
  4. エラーが発生
  ```

- **期待される動作と実際の動作**

- **最小再現コード**
  ```go
  package main
  
  import "github.com/your-username/vibe-coding-logger/pkg/logger"
  
  func main() {
      log := logger.Default()
      // バグが再現するコード
  }
  ```

## ✨ 機能提案

### 提案に含める情報

- **問題の説明**: 現在の制限や課題
- **提案する解決策**: 具体的な機能案
- **使用例**: どのように使われるか
- **実装の考慮点**: 技術的な検討事項
- **代替案**: 他の解決方法

### テンプレート

```markdown
## 問題
現在、○○ができない

## 提案
××機能を追加

## 使用例
\```go
// 使用例のコード
\```

## 実装案
- API設計
- 内部実装
- テスト方針
```

## 📖 ドキュメント貢献

### 改善対象

- README.mdの明確化
- コード例の追加
- APIドキュメントの充実
- 多言語対応

### ドキュメント作成指針

- **明確性**: 初心者でも理解できる
- **完全性**: 必要な情報をすべて含む
- **正確性**: 最新のコードと一致
- **実用性**: 実際の使用場面を想定

## 🌟 ベストプラクティス

### コード品質

- **シンプル**: 複雑さを避ける
- **テスト可能**: テストしやすい設計
- **拡張可能**: 将来の機能拡張を考慮
- **パフォーマンス**: 効率的な実装

### バイブコーディング支援

- **思考プロセス**: 実装の意図を明確に
- **決定記録**: なぜその実装を選択したか
- **学習記録**: 実装中に学んだこと
- **改善点**: 今後の改善案

## 🏆 コントリビューター認定

### レベル

1. **コントリビューター**: 初回貢献
2. **アクティブコントリビューター**: 継続的な貢献
3. **メンテナー**: 重要な機能開発
4. **コアメンテナー**: プロジェクト運営

### 特典

- READMEでの紹介
- 特別なバッジ
- 意思決定への参加権
- 新機能の優先レビュー

## 📞 質問・サポート

### コミュニケーション

- **GitHub Issues**: バグ報告・機能提案
- **GitHub Discussions**: 一般的な質問・議論
- **コードレビュー**: プルリクエストでのフィードバック

### レスポンス時間

- **Issues**: 2-3営業日
- **プルリクエスト**: 1週間以内
- **Discussions**: ベストエフォート

---

## 🙏 最後に

あなたの貢献がVibe Coding Loggerをより良いプロジェクトにします。どんな小さな貢献でも大歓迎です！

**Happy Coding! 🚀**