import '../assets/stylesheets/base.scss';
import React, { Component } from 'react';
import axios from "axios"

import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"


class Home extends Component{

	constructor(props){
		super(props);

		this.state = {
			posts:[]
		};
	}

	componentDidMount() {
		this.getPosts = axios.get("http://127.0.0.1:8080/self/popular/0/")
			.then(function(result) {
				this.setState({posts: result.data.posts});
			}.bind(this))
			.catch(function(e) {
				console.log(e);
			}.bind(this))
	}

	componentWillUnmount() {
		this.getPosts.abort()
	}

	render() {
		return(

			<div>
				{this.state.posts.map((post, index) => {
					return (
						<Post
							key={post.Id}
							Title={post.Title}
							Source={post.Source}
							Description="Description"
						>
						</Post>
					)
				})}
			</div>
		)
	}
}

export default Home;
