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

module.exports = {
	make_magnet: make_magnet,
	chunk: chunk,
	trim: trim
}
