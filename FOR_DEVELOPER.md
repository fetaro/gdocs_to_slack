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


