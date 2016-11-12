"use babel";

import React, { Component } from 'react';
import request from "superagent"

import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"

class SearchResults extends Component{

	constructor(props){
		super(props);

		var posts = [];

		for(var i = 0; i < this.props.posts.length; i++)
		{
			if (this.props.posts[i].error) continue
			posts = posts.concat(this.props.posts[i].body.value);
		}

		posts.sort((a, b) => {
			var aScore = (a.Seeders * 1.1) + a.Leechers;
			var bScore = (b.Seeders * 1.1) + b.Leechers;

			return bScore - aScore;
		});

		this.state = {
			posts: posts
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
