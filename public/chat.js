"use strict";

var socket = new WebSocket('ws://localhost:8080/ws');

// 接続時
socket.onopen = function(event) {
    console.log("接続したよ！");

    // データの用意と送信
    let get = {"method": "getMessage"};
    socket.send(JSON.stringify(get));
    console.log(`こんなデータを送ったよ: ${get}`);
};

// データ受け取り時
socket.onmessage = function(event) {
    console.log(`こんなデータを受け取ったよ: ${event.data}`);

    // 受け取ったデータのパース
    let received_data = JSON.parse(event.data);
    const created_at = new Date(received_data.CreatedAt);

    // メッセージをp要素に変換
    const p = document.createElement("p");
    const text = '[@' + received_data.Name + ']' + '[' + created_at.toLocaleString() + ']' + received_data.Message;
    const message = document.createTextNode(text);
    p.appendChild(message);

    // 画面に表示
    document.getElementById("message_area").prepend(p);
};

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