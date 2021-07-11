import { grpc } from '@improbable-eng/grpc-web';
import { UnaryOutput } from '@improbable-eng/grpc-web/dist/typings/unary';

import { ChatMessage, SendMessageRequest, SubscribeRequest } from './api/service_pb';
import { ChatService } from './api/service_pb_service';

const URL: string = 'http://localhost:5200';

function ready(fn: () => void) {
	if (document.readyState != 'loading'){
		fn();
	} else {
		document.addEventListener('DOMContentLoaded', fn);
	}
}

ready(init);

function init() {
	const johnElements = getElements('john');
	const peterElements = getElements('peter');

	handleAddMessage(johnElements.sendButton, johnElements.input, 'john');
	handleAddMessage(peterElements.sendButton, peterElements.input, 'peter');

	handleMessageReceive('john', johnElements.message);
	handleMessageReceive('peter', peterElements.message);
}

function getElements(user: string): any {
	return {
		input: document.getElementById(user + '-input'),
		sendButton: document.getElementById(user + '-send-btn'),
		message: document.getElementById(user + '-messages'),
	}
}

function handleAddMessage(button: any, input: any, user: string): void {
	button.addEventListener('click', () => {
		const message = input.value;
		input.value = '';
		sendMessage(user, message);
	});
}

function handleMessageReceive(user: string, messages: any): void {
	subscribe(user, (message: ChatMessage.AsObject) => {
		const messageEl = document.createElement('div');
		messageEl.innerText = message.from + ': ' + message.contents;
		messages.appendChild(messageEl);
	});
}

function sendMessage(user: string, message: string) {
	const request = new SendMessageRequest();
    request.setUser(user);
	request.setMessage(message);

	grpc.unary(ChatService.SendMessage, {
		request: request,
		host: URL,
		onEnd: (res: UnaryOutput<any>) => {
			if (res.status !== grpc.Code.OK) {
				console.log('Error calling SendMessage:', res.status, res.statusMessage);
				return;
			}

			console.log('SendMessage called');
		},
	});
}

function subscribe(user: string, callback: (m: ChatMessage.AsObject) => void) {
	const request = new SubscribeRequest();
    request.setUser(user);

    grpc.invoke(ChatService.Subscribe, {
		request: request,
		host: URL,
		onMessage: (e: ChatMessage) => {
			const message = e.toObject();		
			console.log('Message received from gRPC stream:', message);
			callback(message);
		},
		onEnd: () => {
			console.log('Stream ended');
		},
    });
}