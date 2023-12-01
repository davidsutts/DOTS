window.onload = init;

function init() {
	let buttons = document.querySelectorAll("button");
	for (const button of buttons) {
		button.addEventListener('click', () => sendButtonData(button.innerText));
	}
}

function sendButtonData(data: string) {
	let xhr = new XMLHttpRequest;
	xhr.open("POST", "/api/write/" + data);
	xhr.send();
}