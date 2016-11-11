var fs = require("fs");

function make_magnet(infohash) 
{
	// TODO: Include name, trackers, etc. This will work for now though :)
	return "magnet:?xt=urn:btih:" + infohash;
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

module.exports = {
	make_magnet: make_magnet,
	chunk: chunk,
	trim: trim,
	bytes_to_size: bytes_to_size,
	sort: {
		alphanum: alphanum
	},
	loadConfig: loadConfig,
	saveConfig: saveConfig
}
