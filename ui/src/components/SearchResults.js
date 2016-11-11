"use babel";

import React, { Component } from 'react';
import request from "superagent"

import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"

class SearchResults extends Component{

	constructor(props){
		super(props);

		this.state = {
			posts: this.props.posts
		};
	}

	static get defaultProps()
	{ 
		return{
			posts: []
		} 
	}

	render() {
		return(

			<div>
				<h3>Search Results</h3>
				{this.state.posts.map((post, index) => {
					return (
						<Post
							key={post.Id}
							Title={post.Title}
							Source={post.Source}
							Description="Description"
							InfoHash={post.InfoHash}
							Seeders={post.Seeders}
							Leechers={post.Leechers}
						>
						</Post>
					)
				})}
			</div>
		)
	}
}

export default SearchResults;
