# isucrud

ISUCON用のDBへのCRUDを可視化することで、キャッシュすべき箇所の判断を助けるツールです。
リクエストを受け取った関数からどの関数を通ってSQLが実行されているか解析し、以下のようなグラフ構造で表現します。

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
