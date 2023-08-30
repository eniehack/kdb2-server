# kdb2-server

https://kdb2.tsukuba.one のソースコード

## 機能

このソフトウェアはElasticsearchに保管されているシラバスを検索できます。
さらに、シラバスを検索、参照するためのREST APIを提供します。

REST APIに関しては[ここ](https://kdb2.tsukuba.one/api/v0/openapi)にJSON形式のopenapi定義ファイル(v3)があります。

## install

### 依存

- [kdb2-crawler](https://github.com/until-tsukuba/kdb2-crawler)
- [go](https://golang.org)

### setup

1. `git clone` & `cd`
2. kdb2-crawlerを動かし、KdBをクロールし、Elasticsearchにデータを投入する
3. `go build`