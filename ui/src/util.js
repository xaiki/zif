
var fs = require("fs");

function make_magnet(infohash) 
{
	var trackers = [
		"udp://tracker.internetwarriors.net:1337/announce",
		"udp://tracker.leechers-paradise.org:6969/announce",
		"udp://tracker.coppersurfer.tk:6969/announce",
		"udp://exodus.desync.com:6969/announce",
		"udp://tracker.openbittorrent.com:80/announce",
		"udp://tracker.sktorrent.net:6969/announce",
		"udp://tracker.zer0day.to:1337/announce",
		"udp://tracker.pirateparty.gr:6969/announce"
	];

	var link = "magnet:?xt=urn:btih:" + infohash;

	for (var i = 0; i < trackers.length; i++) {
		link += "&tr=" + encodeURIComponent(trackers[i]);
	}

	return link;
}

function chunk(data, size) 
{
	var chunked = [];

	for(var i = 0; i < data.length; i += size)
	{
		var littleChunk = [];

		for(var j = 0; (j < size && i + j < data.length); j++) 
		{
			littleChunk.push(data[i + j]);
		}

		chunked.push(littleChunk);
	}

	return chunked;
}

function trim(string, size, ellipsis) 
{
	var trimmed = string.substring(0, size);
	var left = string.substring(size, string.length);

	if (left.length > 0 && ellipsis)
	{
		trimmed += "...";
		left = "..." + left;
	}

	return {trimmed: trimmed, left: left}
}

// From: http://web.archive.org/web/20130826203933/http://my.opera.com/GreyWyvern/blog/show.dml/1671288
function alphanum(a, b) 
{
	function chunkify(t) 
	{
		var tz = [], x = 0, y = -1, n = 0, i, j;

		while (i = (j = t.charAt(x++)).charCodeAt(0)) 
		{
			var m = (i == 46 || (i >=48 && i <= 57));
			if (m !== n) 
			{
				tz[++y] = "";
				n = m;
			}
			tz[y] += j;
		}
		return tz;
	}

	var aa = chunkify(a);
	var bb = chunkify(b);

	for (var x = 0; aa[x] && bb[x]; x++) 
	{
		if (aa[x] !== bb[x]) 
		{
			var c = Number(aa[x]), d = Number(bb[x]);

			if (c == aa[x] && d == bb[x]) 
			{
				return c - d;

			} else return (aa[x] > bb[x]) ? 1 : -1;
		}
	}

	return aa.length - bb.length;
}

function bytes_to_size(bytes) 
{
   var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

   if (bytes == 0) 
	   return '0 Bytes';

   var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));

   return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
}

function loadConfig()
{
	var conf = fs.readFileSync("./config.json", "utf-8");
	return JSON.parse(conf);
}

function saveConfig(obj)
{
	var json = JSON.stringify(obj);

	fs.writeFile("./config.json", json, (err) => {console.log(err);});
}

function uniq(a) {
    var prims = {"boolean":{}, "number":{}, "string":{}}, objs = [];

    return a.filter(function(item) {
        var type = typeof item;
        if(type in prims)
            return prims[type].hasOwnProperty(item) ? false : (prims[type][item] = true);
        else
            return objs.indexOf(item) >= 0 ? false : objs.push(item);
    });
}

function throttle(fn, threshhold, scope) {
	threshhold || (threshhold = 250);
	var last,
	    deferTimer;
	return function () {
		var context = scope || this;

		var now = +new Date,
		args = arguments;

		if (last && now < last + threshhold) {
			// hold on to it
			clearTimeout(deferTimer);
			deferTimer = setTimeout(function () {
				last = now;
				fn.apply(context, args);
			}, threshhold);

		} else {
			last = now;
			fn.apply(context, args);
		}
	};
}

module.exports = {
	make_magnet: make_magnet,
	chunk: chunk,
	trim: trim,
	bytes_to_size: bytes_to_size,
	sort: {
		alphanum: alphanum
	},
	loadConfig: loadConfig,
	saveConfig: saveConfig,
	uniq: uniq,
	throttle: throttle
}
