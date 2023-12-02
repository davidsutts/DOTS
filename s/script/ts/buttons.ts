window.onload = init;

function init() {
	let buttons = document.querySelectorAll("button");
	for (const button of buttons) {
		button.addEventListener('click', () => {
			sendButtonData(button.innerText);
			updateResults();
		});
	}
}

function sendButtonData(data: string) {
	let xhr = new XMLHttpRequest;
	xhr.open("POST", "/api/write/1/" + data);
	xhr.send();
}

function updateResults() {
	let tableData = document.querySelectorAll("td");

	let xhr = new XMLHttpRequest;
	xhr.open("POST", "/api/read/1");
	xhr.onload = () => {
		let data = JSON.parse(xhr.response);
		for (let i = 1; i <= 5; i++) {
			let objN = data.find((obj) => obj.Value == i);
			if (objN) {
				console.log(objN);
				tableData[i - 1].innerText = objN.Count;
			} else {
				console.log("FAILED");
			}
		}
	};
	xhr.send();

}