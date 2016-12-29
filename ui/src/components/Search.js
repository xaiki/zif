import React, { Component } from 'react';
import request from "superagent"
import async from "async";
import util from "../util"
import TextField from 'material-ui/TextField';
import AutoComplete from 'material-ui/AutoComplete';

class Search extends Component
{
	constructor(props)
	{
		super(props);

		this.state = {
			focus: false,
			searchValue: "",
			focuedWidth: this.props.focusedWidth,
			dataSource: [],
			width: window.innerWidth,
			height: window.innerHeight
		};


		this.onUpdateInput= this.onUpdateInput.bind(this);
		this.onSubmit = this.onSubmit.bind(this);
		this.onResize = this.onResize.bind(this);

		// How many search results we need before they are displayed (and sorted)
		this.searchTotal = 1 + window.config.subscriptions.length;

		this.results = [];

		this.lastEntered = "";
		console.log(this.props.query)

	}

	static get defaultProps() {
		return {
			onResults: function(res){},
			query: "query",
			placeholder: "Search"
		}
	}
	
	onUpdateInput(e)
	{
		if (e.length < 3 || e.length < this.lastEntered.length)
		{ 
			this.lastEntered = e;
			return;
		}

		console.log(e)

		// get completions
		request.post("http://127.0.0.1:8080/self/suggest/")
			.type("form")
			.send({query: e})
			.end((err, res) => {
				this.setState({ dataSource: util.uniq(res.body.value) });
			});
	}

	onSubmit(req)
	{
		var functions = [];
		var subs = window.config.subscriptions;

		// Append local search
		functions.push((cb) => {
			
			request.post("http://127.0.0.1:8080/self/search/")
					.type("form")
					.send({ query: req, page:0 })
					.end(cb)
		});

		for (var i = 0; i < subs.length; i++) 
		{
			var fn = ((i) => { 
				return ((cb) => {
					request.post("http://127.0.0.1:8080/peer/" + subs[i] + "/search/")
							.type("form")
							.send({ query: req, page:0 })

							// Pass null for the error.
							// Weird, I know...
							// This means that even if a peer cannot be connected
							// to, then we still get results from the others.
							// Otherwise an error stops any more results.
							.end((err, res)=>cb(null, res));
				}).bind(this)
			})(i);

			functions.push(fn);
		}

		async.parallel(functions, ((err, res) => {
			if (err) return console.log(err);

			var posts = [];

			for (var i = 0; i < res.length; i++) {
				if (res[i].body) {
					posts.push(res[i].body.value);
				}
			}

			this.props.onResults(posts, req);

		}).bind(this));

	}

	componentWillUnmount() {
	}

	componentDidMount() {
		var search = document.getElementById("search");

		search.addEventListener("keyup", function(event) {
			event.preventDefault();
			if (event.keyCode == 13) {
				this.onSubmit(search.value);
			}
		}.bind(this));
	}

	onResize() {
		this.setState({width: window.innerWidth, height: window.innerHeight})
	}

	static get defaultProps() {
		return {
			placeholder: "Search"
		}
	}
	
	render()
	{
		return (
				<div className="searchBox">
					<div className="searchContainer">
						<span className="icon">
							<i className="material-icons">search</i>
						</span>
						<input type="search" 
							id="search" 
							defaultValue={this.props.query}
							placeholder="Search"
						/>
					</div>
				</div>
		)
	}
}

export default Search;
