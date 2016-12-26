import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';
import request from "superagent"

import FontIcon from 'material-ui/FontIcon';

class NavBar extends Component{

	render(){
		return(
				<ul className="topnav" id="mainMenu">
					<li><span id="logo">Zif</span></li>
					<li><Link to={"/"}>Home</Link></li>
				</ul>
			  )
	}
}

export default NavBar;
