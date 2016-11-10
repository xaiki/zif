import React, { Component } from 'react';
import request from "superagent"

import TextField from 'material-ui/TextField';

class Search extends Component
{
	constructor(props)
	{
		super(props);

		this.state = {
			focus: false,
			searchValue: "",
			focuedWidth: this.props.focusedWidth
		};


		this.toggleFocus = this.toggleFocus.bind(this);
		this.onChange = this.onChange.bind(this);
		this.onSubmit = this.onSubmit.bind(this);
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

	toggleFocus(){ this.setState({ focus: !this.state.focus })} 

	onChange(e)
	{
		this.setState({ searchValue: e.target.value });
	}

	onSubmit(e)
	{
		console.log("Searching for", this.state.searchValue);

		request.post("http://127.0.0.1:8080/self/search/")
				.type("form")
				.send({ query: this.state.searchValue, page:0 })
				.end((err, res) => {
					console.log(res.body)
					this.props.onResults(res.body);
				});

		// stops the page refreshing
		e.preventDefault();
	}

	componentWillUnmount() {
		this.searchRequest.abort();
	}
	
	render()
	{
		this.state.focusedWidth = window.innerWidth - 444;

		var style = {
			backgroundColor: "white",
			paddingLeft: "5px",
			paddingRight: "5px",
			marginRight: "210px",
			marginTop: "8px",
			transition: this.props.transitionTime,
			width: this.state.focus && this.props.growOnFocus ? 
						this.state.focusedWidth : this.props.unfocusedWidth
		};

		return (
			<form onSubmit={this.onSubmit}>
				<TextField
					style={style}
					onFocus={this.toggleFocus}
					onBlur={this.toggleFocus}
					underlineShow={this.props.underlineShow}
					hintText={this.props.hintText}
					onChange={this.onChange}
				/>
			</form>
		)
	}
}

export default Search;
