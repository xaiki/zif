// The Zif HTTP API in JavaScript :)

var request = require("request")

// returns an API object
function setup(url, port) 
{
	var zif = { address: "http://" + url + ":" + port + "/" }

	function make_route (args) 
	{
		var ret = zif.address;

		for (var i = 0; i < args.length; i++)
		{
			ret += args[i] + "/";
		}

		return ret;
	};

	function make_request (args, cb) 
	{
		var route = make_route(args);

		request(route, (err, resp, body) => {
			if (!err && resp.statusCode == 200)
			{
				var data = JSON.parse(body);
				cb(data);
			}
			else if (err)
			{
				console.log(err);
			}
			else if (resp.statusCode != 200) 
			{
				console.log("Error, status code " + resp.statusCode);
			}
		});
	};

	zif.bootstrap = (addr, cb) => make_request(["self", "bootstrap", addr], cb);
	zif.resolve = (addr, cb) => make_request(["self", "resolve", addr], cb);

	zif.recent_remote = (addr, page, cb) => make_request(["peer", addr, "recent", 
			page], cb);
	zif.search_remote = (addr, query, page, cb) => make_request(["peer", addr, "search",
			query], cb);

	zif.search = (query, page, cb) => {
		var route = make_route(["self", "search"])

		request.post(route, { formData: {query: query, page: page} }, (err, resp, body) => {
			if (!err && resp.statusCode == 200)
			{
				var data = JSON.parse(body);
				cb(data);
			}
			else if (err)
			{
				console.log(err);
			}
			else if (resp.statusCode != 200) 
			{
				console.log("Error, status code " + resp.statusCode);
			}
		})
	}

	zif.recent = (page, cb) => make_request(["self", "recent", page], cb);

	return zif;
}

module.exports = setup;
