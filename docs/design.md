# 仕様/技術要件

## 機能要件

1) リアルタイムにチャットを送れる・見れる
2) ユーザーネームはクライアント側に表示、クライアントが決定する　デフォルト：858585
3) ~~ipごとでユーザーネームに色を付ける~~ 保留
4) サーバー側にipバン機能がある 海外ipは弾く
5) ログは最新～30件程度閲覧できる
6) DBには表示分以外のログを残さない
7) ワンページ（画面はチャットルームのみ）
8) 名前、メッセージには最大文字数に制限がある
9) 文字数制限を超える場合、はみ出た部分はすべて1文字の"…"に変換されて正常に書き込まれる

```
[@ユーザーネーム][yyyy/mm/dd hh:mm:ss]おはようございます！
 ↑@●●の部分はipごとに色がつく
```

## 非機能要件

エラー対応、技術制約、保証など

### やること

1) メッセージの時刻はサーバーへの到達順で管理
2) メッセージの送信失敗時はクライアント側にエラーメッセージを表示する

### やらないこと

1) クライアント側の表示順は送信順で保証しない

## 技術選定

1) バックエンド：GoでAPIを用意する
    - DB操作：GORM(仮)
    - 対クライアント：Gin + WebSocket
2) サーバー・クライアント間通信：WebSocketで完結させる
3) フロントエンド：JSでAPIを叩く
4) 動作デバイスはrhaspberry-pi3b上でdocker環境を想定
5) alpine系でgoの公式イメージを想定

## ディレクトリ構成

- app   : バックエンド  API実装まで
    - api: 
    - dbcontroller: dbとのコミュニケーション
- db    : データベース
- public: フロントエンド APIを叩く

### mainパッケージ内各ファイル概略
|ファイル名         |内容概略                       |
|-------------------|-------------------------------|
|main.go            |本体                           |
|db_controller.go   |dbのAPIの実装                  |
|ws_connect.go      |websocketの確立と管理          |
|def_message_type.go|通信に用いるメッセージ型の定義 |
|comm_client.go     |clientとのメッセージのやりとり |
|get_message.go     |getMessageメソッドの実装       |
|post_message.go    |postMessageメソッドの実装      |
|ipban_mw.go        |ip制限をかけるginのミドルウェア|

## API設計

(chat-room-usecase.drawioも参照すること)
- getMessage
    - クライアントが初接続時に叩くAPI　最新のメッセージを規定件数取得する
- postMessage
    - クライアントがメッセージを投稿するときに叩くAPI
- (broadcast)
    - サーバーからメッセージを配信するとき

### getMessage

from client
```json
{
    "method": "getMessage"
}
```

from server
```json
{
    "type": "getReturn",
    "data":
    {
        "status": "OK",
        "count": 10,
        "messages": [
            {
                "name":"ユーザーが決めた名前",
                "message":"投稿するメッセージ(半角1文字以上)",
                "createtime":"2022-10-1 10:10:10.111",
                "isme":false
            },
            {},
        ]
    }
}
```

### postMessage

from client

```json
{
    "method": "postMessage",
    "name":"ユーザーが決めた名前",
    "message":"投稿するメッセージ(半角1文字以上)"
}
```

from server

```json
{
    "type": "postReturn",
    "data":
    {
        "status": "OK"
    }
}
```

```json
{
    "type": "postReturn",
    "data":
    {
        "status": "error"
    }
}
```

### (broadcast)

from server

```json
{
    "type": "broadcast",
    "data":
    {
        "name":"ユーザーが決めた名前",
        "message":"投稿するメッセージ(半角1文字以上)",
        "createtime": "2022-10-1 10:10:10.111",
        "isme": false
    }
}
```

## ダイアグラム

```mermaid
sequenceDiagram
    participant C as Client
    participant OC as OtherClient
    participant S as API
    participant D as DB

    C-->>+S: post Message
    S-->>+D: write DBMessage
    D-->>-S: return OK
    S-->>-C: return OK

    Note right of S: 更新があればメッセージを返す
    loop MessageBroadcast
        S-->>C: return Messages
        S-->>OC: return Messages
    end

    loop DBMessageDelete
        D-->>D: check DBMessage count && delete DBMessage
    end
```

クラス
```mermaid
classDiagram
    class DBMessage{
        +int id
        +String name
        +String message
        +String ipaddress
        +date createtime
        +write(name, message) bool
        +get() bool 
        +delete() bool
    }
```

