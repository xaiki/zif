'use strict';

var electron = require('electron');
var {app, BrowserWindow} = electron;
var spawn = require("child_process").spawn;

let mainWindow;
let torrent;
let hadouken;

function runHadouken() {
	hadouken = spawn("./hadouken", { cwd: "./hadouken" });

	hadouken.stdout.on("data", (data) => {
		console.log("[hadouken]", data.toString());
	});

	hadouken.stderr.on("data", (data) => {
		console.log("[hadouken]", data.toString());
	});
}

function createWindow () {
	mainWindow = new BrowserWindow({width: 800, height: 600});

	mainWindow.loadURL('file://' + __dirname + '/dist/index.html');

	mainWindow.on('closed', function() {
		mainWindow = null;
	});

	console.log("Starting hadouken...")
	runHadouken();
}

app.on('ready', createWindow);

app.on('window-all-closed', function () {
	if (process.platform !== 'darwin') {
		app.quit();
	}
});

app.on('activate', function () {
	if (mainWindow === null) {
		createWindow();
	}
});

