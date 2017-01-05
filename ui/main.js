'use strict';
var electron = require('electron');
var {app, BrowserWindow} = electron;
var TorrentClient = require("./src/TorrentClient.js")


let mainWindow;
let torrent;

function createWindow () {
  mainWindow = new BrowserWindow({width: 800, height: 600});

  mainWindow.loadURL('file://' + __dirname + '/dist/index.html');

  mainWindow.on('closed', function() {
    mainWindow = null;
  });

  torrent = TorrentClient(mainWindow);
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

