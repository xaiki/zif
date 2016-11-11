import React, { Component } from 'react';
import request from "superagent"

import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"

class Home extends Component{

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

	componentDidMount() {
		this.getPosts = request.get("http://127.0.0.1:8080/self/popular/0/")
						.accept("json")
						.type("json")
						.end((err, res) => {
							if (err) {
								return console.log(err);
							}
							this.setState({posts: res.body.posts});
						});

	}

	componentWillUnmount() {
		this.getPosts.abort()
	}

	render() {
		return(

			<div>
				<h3>Popular</h3>
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

export default Home;
