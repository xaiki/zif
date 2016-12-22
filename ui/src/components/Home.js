import React, { Component } from 'react';
import request from "superagent"

import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';
import FlatButton from 'material-ui/FlatButton';

import Post from "./Post"

class Home extends Component{

	constructor(props){
		super(props);

		this.state = {
		};

		if(!this.props.Posts) this.state.posts = [];
		else this.state.posts = this.props.Posts;
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
							this.setState({posts: res.body.value});
						});

	}

	componentWillUnmount() {
		this.getPosts.abort()
	}

	render() {
		console.log(this);
		return(

			<div>
					<h3>Popular</h3>
					<Table>

						<TableHeader>
						  <TableRow>
							<TableHeaderColumn>Name</TableHeaderColumn>
							<TableHeaderColumn>Size</TableHeaderColumn>
							<TableHeaderColumn>Seeders</TableHeaderColumn>
							<TableHeaderColumn>Leechers</TableHeaderColumn>
						  </TableRow>
						</TableHeader>

						<TableBody>
						{this.state.posts.map((post, index) => {
							return (
							<TableRow>
								<TableRowColumn>
								{post.Title}
								</TableRowColumn>
								<TableRowColumn>
								{post.Seeders}
								</TableRowColumn>
								<TableRowColumn>
								{post.Leechers}
								</TableRowColumn>
							</TableRow>
							)
						})}
						</TableBody>
					</Table>
			</div>
		)
	}
}

export default Home;
