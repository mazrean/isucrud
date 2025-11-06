# isucrud

<a href="https://flatt.tech/oss/gmo/trampoline" target="_blank"><img src="https://flatt.tech/assets/images/badges/gmo-oss.svg" height="24px"/></a>

ISUCON用のDBの各テーブルへのCRUDを可視化することで、キャッシュすべき箇所の判断を助けるツールです。
どの関数を通り、どのテーブルへCRUDが発生するSQLが実行されているか解析し、以下のようなグラフ構造で表現します。

[![](./docs/images/private-isu-sample.png)](./docs/sample/private-isu.md)

## Quick Start
1. Goのソースコードのプロジェクトルートに移動
2. isucrudを実行
    ```sh
    isucrud -web ./...
    ```
3. ブラウザでhttp://localhost:7070 にアクセスすると、グラフを見られます

## Install
```sh
go install github.com/mazrean/isucrud@latest
```

## Usage
```
Usage of isucrud:
  -web
      run as web server
  -addr string
      web server address (default ":7070")
  -base string
      base for serving the web server (default "/")
  -dst string
      destination file (default "./dbdoc.md")
  -ignore value
      ignore function
  -ignoreInitialize
      ignore functions with 'initialize' in the name (default true)
  -ignoreMain
      ignore main function (default true)
  -ignorePrefix value
      ignore function
  -version
      show version
```

## 注意事項
SQLの解析を簡易的に行っている関係で、解析に失敗したりCRUD対象テーブルを見逃すことがあります。
あくまでもキャッシュ対象テーブルの判断の参考程度に用いることをお勧めします。
