# isucrud

ISUCON用のDBの各テーブルへのCRUDを可視化することで、キャッシュすべき箇所の判断を助けるツールです。
どの関数を通り、どのテーブルへCRUDが発生するSQLが実行されているか解析し、以下のようなグラフ構造で表現します。

[![](./docs/images/private-isu-sample.png)](./docs/sample/private-isu.md)

## Quick Start
1. Goのソースコードのプロジェクトルートに移動
2. isucrudを実行
    ```sh
    isucrud ./...
    ```
3. `dbdoc.md`をMermaidに対応したマークダウンで開くとグラフを見れます

## Install
```sh
go install github.com/mazrean/isucrud@latest
```

## 注意事項
SQLの解析を簡易的に行っている関係で、解析に失敗したりCRUD対象テーブルを見逃すことがあります。
あくまでもキャッシュ対象テーブルの判断の参考程度に用いることをお勧めします。
