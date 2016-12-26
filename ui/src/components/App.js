import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';

import {Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle} from 'material-ui/Toolbar';
import {grey100, grey50} from 'material-ui/styles/colors';


import Home from './Home';
import SearchResults from "./SearchResults"
import Stream from "./Stream"
import Welcome from "./WelcomeDialog"

import util from "../util"

var WebTorrent = require("webtorrent");

var routes = [{ path: "/", component: Home },
			  { path: "/search", component: SearchResults },
			  { path: "/stream/:infohash", component: Stream },
			  { path: "/downloads", component: Stream }];

class App extends Component
{

	constructor(props)
	{
		super(props);

		this.state = { drawerOpen : true, 
			search: {
				focus: false
			}
		};

		this.handleToggle = this.handleToggle.bind(this);
		this.onResults = this.onResults.bind(this);

		this.config = util.loadConfig();

		window.downloadClient = new WebTorrent();

		window.zifColor = {
			primary: "#3f3b3b",
			secondary: "#eee9d9",
			highlight: "#DE1B1B",
			accent: "#E9E581"
		};
	}

	handleToggle(){ this.setState({ drawerOpen: !this.state.drawerOpen }) }

	onResults(res) 
	{
		routes[1].component = () => {
			return (
				<SearchResults posts={res}/>
			)
		};

		hashHistory.push("/search")
	}

	homeButtonClick(){ hashHistory.push("/") }

	render() 
	{
		var style = {
			drawerItems: {
				marginTop: "75px"
			},

			router: {
				marginTop: "10px"
			},
			topbar: {
				backgroundColor: window.zifColor.primary
			}

		}

		return(
			<div style={{height: "100%"}}>
				<Welcome config={this.config}/>


				<div>
					<Router history={hashHistory} routes={routes}>
					</Router>
				</div>
			</div>
		)
	}
}

export default App;
