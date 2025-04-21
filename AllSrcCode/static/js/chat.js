let ws;
document.addEventListener('DOMContentLoaded', function() {
    // 连接 WebSocket
    connectWebSocket();

    // 获取用户列表
    // getUsers();

    // 监听消息输入
    document.getElementById('messageInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });
});

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    // const url = protocol + window.location.host + '/ws?userID=' + currentUser;
    const url = protocol + window.location.host + '/ws';
    ws = new WebSocket(url);

    ws.onopen = function() {
        console.log('Connected to WebSocket server');
        getUsers();
        // 设置定时器，每3秒更新一次用户列表
        userListUpdateInterval = setInterval(getUsers, 3000);
    };

    ws.onmessage = function(event) {
        const message = event.data;
        appendMessage(message);
    };

    ws.onclose = function() {
        console.log('Disconnected from WebSocket server');
        clearInterval(userListUpdateInterval); // 清除定时器
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

function getUsers() {
    fetch('/users')
        .then(response => response.json())
        .then(users => {
            const userList = document.getElementById('userList');
            const userSelect = document.getElementById('userSelect');

            userList.innerHTML = '';
            userSelect.innerHTML = '<option value="all">广播</option>';

            users.forEach(user => {
                if (user !== currentUser) {
                    const li = document.createElement('li');
                    li.textContent = user;
                    userList.appendChild(li);

                    const option = document.createElement('option');
                    option.value = user;
                    option.textContent = user;
                    userSelect.appendChild(option);
                }
            });
        });
}

function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const userSelect = document.getElementById('userSelect');
    const message = messageInput.value;

    if (message.trim() === '') return;

    const messageObj = {
        receiver_id: userSelect.value,
        text: message
    };

    ws.send(JSON.stringify(messageObj));
    messageInput.value = '';
}

function appendMessage(message) {
    const messagesDiv = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.textContent = message;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}