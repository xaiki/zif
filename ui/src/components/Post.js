import React, { Component } from 'react';
import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import axios from "axios"

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
				<Card>
					<CardHeader
						title={this.props.Title}
						subtitle={this.props.Source}
						actAsExpander={true}
						showExpandableButton={true}
					/>
					<CardText expandable={true}>
						{this.props.Description}
					</CardText>

					<CardActions>
						<FlatButton label="Download" />
						<FlatButton label="Stream" />
					</CardActions>
				</Card>
			</div>)
	}
}

export default Post;
