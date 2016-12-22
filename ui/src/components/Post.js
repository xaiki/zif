import React, { Component } from 'react';
import Chip from 'material-ui/Chip';
import { hashHistory, Link } from 'react-router';
import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import {List, ListItem} from 'material-ui/List';

const style = {
	marginTop: "10px"
};

class Post extends Component
{
	constructor(props){
		super(props);
	}

	render() {
		return (
			<div style={style}>
				<ListItem primaryText={this.props.Title} />
			</div>)
	}
}

export default Post;
