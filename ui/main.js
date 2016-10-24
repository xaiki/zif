const {app, BrowserWindow} = require("electron")
const spawn = require("child_process").spawn;

let win
let zifd

function createWindow() 
{
	win = new BrowserWindow({ width: 800, height: 600 });

	win.loadURL(`file://${__dirname}/index.html`);
	//win.webContents.openDevTools();

	win.on("closed", () => {
		win = null;
	})

	// TODO: Make this optional. Some users may well be running a remote daemon,
	// or may have one running anyway in order to use other clients. Who knows?
	zifd = spawn("zifd", [], {stdio: "inherit"});
}

app.on("ready", createWindow);
