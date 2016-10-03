const { app, BrowserWindow } = require("electron")

// to avoid the window being closed by GC
let win

function createWindow() 
{
	win = new BrowserWindow({ width: 800, height: 600 });

	win.loadURL(`file://${__dirname}/html/index.html`)
	win.webContents.openDevTools();

	win.on("closed", () => {

		win = null;

	});
}

// app events

app.on("ready", createWindow);

app.on("window-all-closed", () => {
	app.quit()
});
