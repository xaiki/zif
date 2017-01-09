var request = require("superagent");

var hadouken = {};

function parseTorrent(arr) {
	var torrent = {
		infohash: arr[0],
		status: arr[1],
		name: arr[2],
		size: arr[3],
		progress: arr[4],
		downloaded: arr[5],
		uploaded: arr[6],
		ratio: arr[6],
		up: arr[8],
		down: arr[9],
		eta: arr[10],
		label: arr[11],
		numLeech: arr[12],
		listPeers: arr[13],
		numSeeds: arr[14],
		listSeeds: arr[15],
		distributionCopies: arr[16],
		queuePosition: arr[17],
		remaining: arr[18],
		downloadUrl: arr[19],
		rssUrl: arr[20],
		error: arr[21],
		streamId: arr[22],
		addedTime: arr[23],
		completedTime: arr[24],
		updateUrl: arr[25],
		savePath: arr[26]
	};

	return torrent;
}

hadouken.method = function(method, params, cb){
	request.get("localhost:7070/api")
		.send(JSON.stringify({ method: method, params: params }))
		.auth("admin", "admin")
		.end((err, res) => {
			if (cb) cb(err, res);
		});
}

hadouken.list = function(cb){
		this.method("webui.list", [], (err, res) => {
			if (err) {
				cb(err, null);
				return;
			}

			cb(err, res.body.result.torrents.map(parseTorrent))
		});
}

hadouken.addLink = function(link, cb) {
		this.method("webui.addTorrent", ["url", link],(err, res) => {
			if (err) {
				cb(err, null);
				return;
			}

			cb(err, res.body)
		});
}

module.exports = hadouken;
