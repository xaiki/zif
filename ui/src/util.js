function make_magnet(infohash) 
{
	// TODO: Include name, trackers, etc. This will work for now though :)
	return "magnet:?xt=urn:btih:" + infohash;
}

module.exports = {
	make_magnet: make_magnet
}
