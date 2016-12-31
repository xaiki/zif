import React, { Component } from 'react';
import { hashHistory, Link } from 'react-router';
import {Card, CardActions, CardHeader, CardText, CardMedia} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';

const style = {
	marginTop: "10px"
};

class StreamFile extends Component
{
	constructor(props){
		super(props);

		this.state = {
			expanded: false
		};

		this.onExpandChange = this.onExpandChange.bind(this);
	}

	onExpandChange(newState)
	{
		this.setState({expanded: newState});

	}

	componentDidUpdate()
	{
		var target = this.refs["s" + this.props.Index];

		if (!this.state.expanded) return;

		if (this.elem) 
		{
			target.appendChild(this.elem);
			return;
		}

		console.log(this.props.File);

		this.props.File.appendTo(target, 
				(err, elem) => {
					this.elem = elem;
					this.elem.style.width = "100%";
				});
	}

	render() {
		var target = this.refs["s" + this.props.Index];
		console.log(target)

		return (
			<div style={style}>
				<Card onExpandChange={this.onExpandChange}>
					<CardHeader
						title={this.props.Title}
						actAsExpander={true}
						showExpandableButton={true}
					/>
					<CardText expandable={true}>
						<div ref={"s" + this.props.Index}></div>
					</CardText>
				</Card>
			</div>)
	}
}

export default StreamFile;
