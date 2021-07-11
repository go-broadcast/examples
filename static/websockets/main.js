window.onload = function() {
	initUser('john');
	initUser('peter');
}

function initUser(user) {
	var elements = getElements(user);

	var webSocket = new WebSocket('ws://localhost:5200/ws/' + user);

	elements.sendButton.addEventListener('click', function() {
		var message = elements.input.value;
		elements.input.value = '';
		webSocket.send(message);
	});

	webSocket.onclose = function(ev) {
		console.log('socket close:', ev);
	};

	webSocket.onerror = function(ev) {
		console.log('socket error:', ev);
	};

	webSocket.onmessage = function(ev) {
		var messageEl = document.createElement('div');
		messageEl.innerText = JSON.parse(ev.data);
		elements.messages.appendChild(messageEl);

		console.log('socket message:', ev.data);
	};
}

function getElements(user) {
	return {
		input: document.getElementById(user + '-input'),
		sendButton: document.getElementById(user + '-send-btn'),
		messages: document.getElementById(user + '-messages'),
	}
}
