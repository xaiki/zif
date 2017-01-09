const WebTorrent = require("webtorrent-hybrid");
const {ipcMain} = require("electron");
const fs = require("fs");

var request = require("superagent");


function TorrentClient(web){
	var client = {};

	client.web = web;
	client.seed = seed.bind(client);
	client.download = download.bind(client);
	client.client = new WebTorrent();

	client.client.on("error", (err) => {
		console.log(err);
	});

	ipcMain.on("seed", client.seed);
	ipcMain.on("download", client.download);
	ipcMain.on("torrents", () => web.send("torrents", client.client.torrents));

	return client;
}

function download(e, arg) {
	console.log("Downloading ", arg);

	this.client.add(arg, (torrent) => {
		this.web.send("torrent", torrent);
	});
}

function seed(e, arg) {
	var dirRegex = "^(.+)/([^/]+)$";
	var dir = arg.files[0].match(dirRegex)[1];


	console.log("adding torrent")

	this.client.seed(arg.files, {path: dir}, (torrent) => {
		var data = {
			title: arg.title,
			meta: JSON.stringify(arg.meta),
			infoHash: torrent.infoHash
		};

		request.post("http://127.0.0.1:8080/self/addpost/")
			.type("form")
			.send({data: JSON.stringify(data), index: true})
			.end((err, res) => {
			});

		var fileName = "./torrents/" + torrent.name + ".torrent";
		fs.writeFile(fileName, torrent.torrentFile, (err) => console.log(err));

		this.web.send("torrent", torrent);
	});
}

		
module.exports = TorrentClient;
