var torrent = require("webtorrent");
var http = require("http");
var util = require("../util");

function torrentStream(ipc){
	var ret = {};
	ret.ipc = ipc;
	ret.client = new torrent();

	ipc.on("stream-magnet", (e, arg) => {
		if (ret.server)
			ret.server.close();

		var t = ret.client.get(arg);

		if (t) {
			ret.server = t.createServer();
			ret.server.listen(60000);
			e.sender.send("torrent", t);
			return;
		}

		ret.client.add(arg, {path: "$HOME/Downloads"}, (torrent) => {
			console.log("downloading", torrent.infoHash);

			ret.server = torrent.createServer();
			ret.server.listen(60000);

			e.sender.send("torrent", torrent);
		});
	});

	return ret;
}

module.exports = torrentStream;
