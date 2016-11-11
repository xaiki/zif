import React, { Component } from 'react';
import Chip from 'material-ui/Chip';
import { hashHistory, Link } from 'react-router';
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
						showExpandableButton={true}>

						<div style={{display: "flex", flexWrap: "wrap"}}>
							<Chip>{this.props.Seeders} Seeders</Chip>
							<Chip>{this.props.Leechers} Leechers</Chip>
						</div>

					</CardHeader>
					<CardText expandable={true}>
						{this.props.Description}
					</CardText>

					<CardActions>
						<FlatButton label="Download" />

						<Link to={"/stream/" + this.props.InfoHash}>
							<FlatButton label="Stream" />
						</Link>
					</CardActions>
				</Card>
			</div>)
	}
}

export default Post;
