import util from "./util.js"
import WebTorrent from "webtorrent"

// Add more of these.
var streamable = ["mp3", "mp4", "ogg", "webm"];

function add(link) 
{
	console.log(this);
	this.client.add(link, this.added);
}

function added(torrent)
{
	console.log("Downloading:", torrent.infoHash);

	console.log(this);
	this.files = torrent.files;

	// Sort by filename
	this.files.sort(util.sort.alphanum);
}

// Add files to the HTML
function add_files()
{
	this.component.entries = this.files;
}

function stream(component)
{
	var stream = {
		client: new WebTorrent(),
		add: add,
		added: added,
		component: component,
		add_files: add_files
	}

	return stream;
}

module.exports = stream;
