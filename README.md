# 目的

Google Docsの箇条書きをSlackにコピーするツール

## 課題

Google Docsの箇条書きをMacのクリップボードにコピーしてSlackにペーストすると、以下のように箇条書きが崩れてしまう。

例えば、このようなGoogle Docsの箇条書き

<img src="docs/docs.png" width="400" alt="good_slack">

これをSlackにコピーすると

<img src="docs/bad_slack.png" width="400" alt="bad_slack">

とフラットな箇条書きでコピーされてしまう

### このツールを使うと

インデントの構造を保持してコピーできます！

<img src="docs/good_slack.png" width="400" alt="good_slack">

## インストール

リポジトリに含まれるバイナリをそのまま使用できます。

```bash
git clone https://github.com/fetaro/docks_to_slack_go.git
cd docks_to_slack_go
# 必要に応じてパスの通った場所にコピーしてください
cp dist/docs_to_slack /usr/local/bin/
```

## 使い方

1. Google Docsの箇条書きをコピーします

2. クリップボードにGoogle Docsの箇条書きがある状態で以下のコマンドを実行します

```bash
docs_to_slack
```

これでSlackに貼り付け可能なデータがクリップボードに書き込まれました

3. Slackにてペーストします

### (オプション)テキストとしてペーストする場合

`-t` オプションでmarkdown形式のテキストとしてペーストすることも可能です

```bash
docs_to_slack -t
```

以下のようなテキストがクリップボードに保存されます

```
- インデントレベル１
    - インデントレベル２のA
    - インデントレベル２のB
    - インデントレベル２のC
```

### デバッグ

`-d` オプションを付けるとデバッグモードになります

```bash
docs_to_slack -d
```

## 開発者用

[FOR_DEVELOPER.md](FOR_DEVELOPER.md) を参照してください。
