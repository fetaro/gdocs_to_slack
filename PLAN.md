# 実装計画: docs_to_slack_go

## 概要
Python製ツール `docs_to_slack` をGo言語に移植し、macOS上で単一バイナリとして動作するようにします。
クリップボード内のHTML（Google Docs等のリスト）を読み取り、Slack貼り付け用の形式（`slack/texty`）に変換して書き戻します。

## アーキテクチャ
`simplify_clipboard_html` を参考に、CGOを用いてmacOSのネイティブAPI（Cocoa/AppKit）を呼び出します。

### 構成要素
1.  **Main / CLI**: コマンドライン引数の処理、全体のフロー制御。
2.  **Clipboard Manager (CGO)**:
    *   `pbpaste` 相当の HTML 読み込み。
    *   `org.chromium.web-custom-data` 形式へのバイナリ書き込み。
3.  **HTML Parser**: `golang.org/x/net/html` を使用してHTMLを解析し、テキストとJSON構造を生成。
4.  **Pickle Writer**: Chromiumのカスタムデータ形式（Pickle）に従い、バイナリデータを生成。

## 実装ステップ

### Phase 1: プロジェクト初期化とクリップボード読み込み
*   `go.mod` の作成。
*   `simplify_clipboard_html` を参考に、CGOを使ってクリップボードからHTMLを取得する機能を実装。
*   取得したHTMLを標準出力するだけのプロトタイプを作成。

### Phase 2: HTML解析ロジックの移植
*   Python版 `SlackListGenerator` クラスのロジックを移植。
*   `<ul>`, `<ol>`, `<li>` の構造解析と `aria-level` の処理。
*   出力:
    *   プレーンテキスト（インデント付き）
    *   Slack用JSON (`ops` リスト)

### Phase 3: バイナリデータ生成 (Pickle Writer)
*   Python版 `PickleWriter` クラスのロジックを移植。
*   Chromium Pickleフォーマットの実装:
    *   Little Endian
    *   UTF-16LE 文字列エンコーディング
    *   4バイトアライメント
*   `public.utf8-plain-text` と `slack/texty` を含むペイロードの作成。

### Phase 4: クリップボード書き込みと統合
*   CGO側にバイナリデータを書き込む関数 `setDataToPasteboard` を追加。
*   `NSPasteboard` を使い、`org.chromium.web-custom-data` タイプで書き込む。
*   CLIオプション (`-d`, `-t`) の実装。

### Phase 5: ビルドと配布準備
*   `build.sh` の作成（`simplify_clipboard_html` 参考）。
*   動作確認。

## 参考資料
*   Python実装: `../docs_to_slack/`
*   Go実装参考: `../simplify_clipboard_html/`
