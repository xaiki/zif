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
		this.searchTotal = 1 + this.props.Subscriptions.length;

		this.results = [];

		this.lastEntered = "";
	}

	static get defaultProps()
	{ 
		return{
			underlineShow: false,
			hintText: "Search...",
			growOnFocus: true,
			unfocusedWidth: "256px",
			focusedWidth: "512px",
			transitionTime: ".3s",
			onResults: function(){}
		} 
	}

	static get searchWidth()
	{
		return window.innerWidth - 446;
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

	onSubmit(req, i)
	{
		var functions = [];

		// Append local search
		functions.push((cb) => {
			
			request.post("http://127.0.0.1:8080/self/search/")
					.type("form")
					.send({ query: req, page:0 })
					.end(cb)
		});

		for (var i = 0; i < this.props.Subscriptions.length; i++) 
		{
			var fn = ((i) => { 
				return ((cb) => {
					request.post("http://127.0.0.1:8080/peer/" + this.props.Subscriptions[i] + "/search/")
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

		async.parallel(functions, (err, res) => {
			if (err) return console.log(err);
			console.log(res)
			this.props.onResults(res);
		});

	}

	componentWillUnmount() {
		this.searchRequest.abort();
	}

	componentDidMount() {
		window.addEventListener("resize", this.onResize);
	}

	onResize() {
		this.setState({width: window.innerWidth, height: window.innerHeight})
	}
	
	render()
	{
		var style = {
			backgroundColor: "white",
			paddingLeft: "5px",
			paddingRight: "5px",
			marginRight: "150px",
			marginTop: "8px",
			width: this.state.width - 333
		};

		return (
				<div>
			<AutoComplete
				style={style}
				dataSource={this.state.dataSource}
				underlineShow={this.props.underlineShow}
				hintText={this.props.hintText}
				onUpdateInput={util.throttle(this.onUpdateInput, 100, this)}
				filter={AutoComplete.fuzzyFilter}
				fullWidth={true}
				onNewRequest={this.onSubmit}
			/>
			</div>
		)
	}
}

export default Search;
