# go-rpc-server

Go で RPC スタイルの API を実装するサンプル

## 起動方法

```sh
docker build -t go-rpc-server .
docker run --rm -it -p 8000:8000 go-rpc-server
```

## 動作確認コマンド

```sh
curl -X POST localhost:8000/core/v1/task/get -d '{ "id": "id_01" }'
curl -X POST localhost:8000/core/v1/task/list
curl -X POST localhost:8000/core/v1/task/create -d '{ "title": "title_01" }'
curl -X POST localhost:8000/core/v1/task/update -d '{ "id": "id_01", "title": "title_01", "status": "DONE" }'
curl -X POST localhost:8000/core/v1/task/delete -d '{ "id": "id_01" }'
```
