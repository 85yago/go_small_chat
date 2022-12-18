"use strict";

var socket = new WebSocket('ws://localhost:8080/ws');

// 接続時
socket.onopen = function(event) {
    console.log("接続したよ！");

    // データの用意と送信
    let get = {"method": "getMessage"};
    socket.send(JSON.stringify(get));
    console.log(`こんなデータを送ったよ: ${JSON.stringify(get)}`);
};

// メッセージを画面に表示する関数
function addMessage(msg) {
    const createTime = new Date(msg.createtime);

    // メッセージをp要素に変換
    const p = document.createElement("p");

    p.appendChild(document.createTextNode(`[`));

    const user_name = document.createElement("span").appendChild(document.createTextNode('@' + msg.name));
    p.appendChild(user_name);
    
    const text = ']' + '[' + createTime.toLocaleString() + ']' + msg.message;
    const message = document.createTextNode(text);
    p.appendChild(message);

    // 画面に表示
    document.getElementById("message_area").prepend(p);
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

// データ受け取りイベント時に登録
socket.addEventListener("message", getReturn);
socket.addEventListener("message", postReturn);
socket.addEventListener("message", receiveBroadcast);

// 接続終了時
socket.onclose = function(event) {
    console.log('接続が死んだよ…');
};

// エラー時
socket.onerror = function(error) {
    console.log(`接続がエラったよ: ${error.message}`);
};

document.addEventListener('DOMContentLoaded',function(event){
    // ボタン押した時の送信用関数
    document.getElementById('send').addEventListener('click', function(event){
        event.preventDefault();

        // データの準備
        let send_data = {};
        send_data["method"] = "postMessage";
        send_data["name"] = document.getElementById('name').value;
        send_data["message"] = document.getElementById('message').value;

        // データの送信
        socket.send(JSON.stringify(send_data));
    });
});