# 開発者用ドキュメント

## 開発環境のセットアップ

Go 1.24 以上が必要です。

```bash
git clone https://github.com/fetaro/docks_to_slack_go.git
cd docks_to_slack_go
go mod download
```

## ビルド

```bash
./build.sh
```

または

```bash
go build -o docs_to_slack
```

## テスト

Goの標準テストフレームワークを使用しています。

```bash
go test -v .
```

### テスト内容

*   `html_parser_test.go`: HTML解析ロジックのテスト。ネストされたリストや `aria-level` の処理が正しく行われるか検証します。
*   `pickle_writer_test.go`: Chromium Pickleフォーマット（バイナリデータ）の生成ロジックのテスト。リトルエンディアン、UTF-16エンコーディング、パディングなどが正しく処理されるか検証します。

## プロジェクト構造

*   `main.go`: エントリーポイント。CLI引数の処理と、CGOを使用したクリップボード操作（読み込み・書き込み）を行います。
*   `html_parser.go`: HTML文字列を解析し、プレーンテキストとSlack用JSON (`slack/texty`) を生成します。
*   `pickle_writer.go`: Chromiumのカスタムクリップボード形式 (`org.chromium.web-custom-data`) に準拠したバイナリデータを生成します。
*   `PLAN.md`: 実装計画書。

## 参考

*   Python実装: `../docs_to_slack/`
*   Go実装参考: `../simplify_clipboard_html/`
