import React, { Component } from 'react';
import request from "superagent"

import FlatButton from 'material-ui/FlatButton';
import FontIcon from 'material-ui/FontIcon';

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
				<div className="searchBox">
					<div className="searchContainer">
						<span className="icon">
							<FontIcon className="material-icons">search</FontIcon>
						</span>
						<input type="search" id="search" placeholder="Search" />
					</div>
				</div>
			</div>
		)
	}
}

export default Home;
