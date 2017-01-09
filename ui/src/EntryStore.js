import request from "superagent"

var entries = {};

function resolve(address, cb) {
	if (entries[address]){ 
		cb(null, entries[address]);
		return entries[address];
	}

	request.get("http://127.0.0.1:8080/self/resolve/" + address + "/")
		.end((err, res) => {
			if (err) return;

			entries[address]= res.body.value;

			if (cb) cb(err, res.body.value);
		});
}

export default resolve;
