"use strict";

const MAX_DISPLAY_MSG_COUNT = 100;
const TRY_RECONNECT_COUNT = 5;
// デバッグ用にURLを変える関数
function retUrl() {
    if (location.protocol == "https:") {
        return "wss://" + location.host + "/ws";
    } else if (location.protocol == "http:") {
        return "ws://" + location.host + "/ws";
    }
}
const WS_URL = retUrl();
let socket = new WebSocket(WS_URL);

// 接続時
function forOpen(event) {
    console.log("接続したよ！");

    // データの用意と送信
    let get = { "method": "getMessage" };
    socket.send(JSON.stringify(get));
    console.log(`こんなデータを送ったよ: ${JSON.stringify(get)}`);
};

// エラー時
function forError(error) {
    console.log(`接続がエラったよ: ${error.message}`);
};

// 接続終了時
function forClose(event) {
    console.log('接続が死んだよ…');
    let isReconnect = false;

    for (let i = 0; i < TRY_RECONNECT_COUNT; i++) {
        socket = new WebSocket(WS_URL);
        // CONNECTINGは微妙かも
        if (socket.readyState == WebSocket.CONNECTING) {
            isReconnect = true;
            break;
        }
    }

    if (isReconnect) {
        console.log("接続し直したよ！");
        setTimeout(() => {
            registorFunctions();
        }, 500);
    } else {
        throw "socket's error: dead socket.";
    }
};

// メッセージを画面に表示する関数
function addMessage(msg) {
    const createTime = new Date(msg.createtime);

    // メッセージをp要素に変換
    const p = document.createElement("p");

    // 内容を連結する要素に加工
    const userName = document.createElement("span").appendChild(document.createTextNode('[@' + msg.name + ']'));
    const time = document.createTextNode('[' + createTime.toLocaleString() + ']');
    const message = document.createTextNode(msg.message);

    // p要素に追加
    p.appendChild(userName);
    p.appendChild(time);
    p.appendChild(message);

    // 画面に表示
    document.getElementById("message_area").prepend(p);

    // MAX_DISPLAY_MSG_COUNTを上回った時に消す
    while (MAX_DISPLAY_MSG_COUNT <= p.childElementCount) {
        p.removeChild(p.lastChild);
    }
}

// getReturn受け取り用
function getReturn(event) {
    // エラーチェックとか
    let received_data;
    try {
        received_data = JSON.parse(event.data);
    } catch (error) {
        console.log(error);
        return;
    }

    if (received_data.type !== "getReturn") {
        return;
    }
    console.log("getReturn.");

    if (received_data.data.status !== "OK") {
        throw `getReturn's error: ${received_data.data.status}`;
    }

    // メッセージを画面に表示する
    const count = received_data.data.count;
    const msgs = received_data.data.messages;
    for (let i = 0; i < count; i++) {
        addMessage(msgs[i]);
    }
}

// postReturn受け取り用
function postReturn(event) {
    // エラーチェックとか
    let received_data;
    try {
        received_data = JSON.parse(event.data);
    } catch (error) {
        console.log(error);
        return;
    }

    if (received_data.type !== "postReturn") {
        return;
    }
    console.log("postReturn.");

    // status確認
    if (received_data.data.status !== "OK") {
        alert(`postMessage fail!: ${received_data.data.status}`);

        throw `postReturn's error: ${received_data.data.status}`;
    }
}

// broadcast受け取り用
function receiveBroadcast(event) {
    // エラーチェックとか
    let received_data;
    try {
        received_data = JSON.parse(event.data);
    } catch (error) {
        console.log(error);
        return;
    }

    if (received_data.type !== "broadcast") {
        return;
    }
    console.log("receiveBroadcast.");

    // メッセージを画面に表示する
    addMessage(received_data.data);
}

// 登録関数
function registorFunctions() {
    socket.addEventListener("message", getReturn);
    socket.addEventListener("message", postReturn);
    socket.addEventListener("message", receiveBroadcast);
    socket.addEventListener("open", forOpen);
    socket.addEventListener("close", forClose);
    socket.addEventListener("error", forError);
}

// ここから実行される

document.addEventListener('DOMContentLoaded', function (event) {
    // ボタン押した時の送信用関数
    document.getElementById('send').addEventListener('click', function (event) {
        event.preventDefault();

        // データの準備
        let send_data = {};
        send_data["method"] = "postMessage";
        send_data["name"] = document.getElementById('input_name').value;
        send_data["message"] = document.getElementById('input_message').value;

        // データの送信
        socket.send(JSON.stringify(send_data));

        // メッセージ欄を空にする
        document.getElementById('input_message').value = "";
    });
});

registorFunctions();
