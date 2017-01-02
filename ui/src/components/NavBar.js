import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';
import request from "superagent"

import util from "../util.js"
import ReactTooltip from 'react-tooltip'

class NavBar extends Component{

	render(){
		return(
				<ul className="topnav" id="mainMenu">
					<li><span id="logo">Zif</span></li>
					<li><Link className="navitem" to={"/"}>Home</Link></li>
					<li style={{float: "right", height:0}}>
						<span>
							<a data-tip id="upload" className="navitem">
								<i className="material-icons">file_upload</i>
							</a>
						</span>
					</li>

					<ReactTooltip data-for="upload" place="left" effect="solid">
						<span>Upload</span>
					</ReactTooltip>
				</ul>

			  )
	}
}

export default NavBar;
