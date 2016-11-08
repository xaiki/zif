import React, { Component } from 'react';
import axios from "axios"

import TextField from 'material-ui/TextField';

class Search extends Component
{
	constructor(props)
	{
		super(props);

		this.state = {
			focus: false,
			searchValue: ""
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
			unfocusedWidth: "256px",
			focusedWidth: "512px",
			growOnFocus: true,
			transitionTime: ".3s"
		} 
	}

	toggleFocus()
	{
		this.setState({ focus: !this.state.focus });
	}

	onChange(e)
	{
		this.setState({ searchValue: e.target.value });
	}

	onSubmit(e)
	{
		console.log("Searching for", this.state.searchValue);
	}

	componentWillUnmount() {
		this.searchRequest.abort();
	}
	
	render()
	{
		var style = {
			backgroundColor: "white",
			paddingLeft: "5px",
			paddingRight: "5px",
			marginRight: "210px",
			marginTop: "8px",
			transition: this.props.transitionTime,
			width: this.state.focus && this.props.growOnFocus ? 
						this.props.focusedWidth : this.props.unfocusedWidth
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
