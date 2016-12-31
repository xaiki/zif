"use babel";

import React, { Component } from 'react';
import request from "superagent"

import {List, ListItem} from 'material-ui/List';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"
import NavBar from "./NavBar"
import Search from "./Search"

class SearchResults extends Component{

	constructor(props){
		super(props);


		this.state = {};
		this.state.posts = this.sortPosts(this.props.location.state.posts);
	}

	static get defaultProps()
	{ 
		return{
			posts: [],
			value: ""
		} 
	}

	sortPosts(posts) {
		var newPosts = [];

		for(var i = 0; i < posts.length; i++)
		{
			if (!posts[i]) continue;

			for (var j = 0; j < posts[i].length; j++) {
				newPosts.push(posts[i][j]);
			}
		}

		newPosts.sort((a, b) => {
			var aScore = (a.Seeders * 1.1) + a.Leechers;
			var bScore = (b.Seeders * 1.1) + b.Leechers;

			return bScore - aScore;
		});

		return newPosts;
	}

	render() {
		return(

			<div className="content">
				<NavBar />
				<Search 
					query={this.props.location.state.query} 
					router={this.props.router}
					onResults={(res, query) => {
						this.setState({
							posts: this.sortPosts(res)
						})
					}}/>

				<div className="searchResults">
					{this.state.posts.map((post, index) => {
						return (
							<Post key={index}
								  title={post.Title}
								  infohash={post.InfoHash}
								  seeders={post.Seeders}
								  leechers={post.Leechers}
								  meta={post.Meta}
								  size={post.Size}
								  fileCount={post.FileCount}
								  tags={post.Tags}
								  uploadDate={post.UploadDate}/>
						)
					})}
				</div>
			</div>
		)
	}
}

export default SearchResults;
