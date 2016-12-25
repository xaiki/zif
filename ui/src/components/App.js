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
			drawer: {
				backgroundColor: grey100,
				zIndex: 1000
			},

			drawerItems: {
				marginTop: "75px"
			},

			router: {
				paddingLeft: "210px",
				paddingRight: "210px",
				paddingBottom: "10px",
				marginTop: "75px"
			},
			topbar: {
				backgroundColor: "red"
			}

		}

		return(
			<div style={{height: "100%"}}>
				<Welcome config={this.config}/>

				<Toolbar>
					<ToolbarGroup firstChild={true}>
						<ToolbarTitle text="Zif" style={ {marginLeft:"10px"}}/>
					</ToolbarGroup>
				</Toolbar>

				<div style={style.router}>
					<Router history={hashHistory} routes={routes}>
					</Router>
				</div>
			</div>
		)
	}
}

export default App;
