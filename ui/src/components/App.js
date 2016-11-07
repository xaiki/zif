import '../assets/stylesheets/base.scss';
import '../assets/stylesheets/app.scss';

import React, { Component } from 'react';
import { Router, Route, hashHistory } from 'react-router'

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';

import Home from './Home';


class App extends Component{
  render() {
    return(
    	<div>
    		<AppBar title="Zif"/>

			<Drawer width={200} docked={true} open={true} className="drawer">
				<MenuItem>Home</MenuItem>
			</Drawer>

			<div id="router">
				<Router history={hashHistory}>
					<Route path="/" component={Home}/>
				</Router>
			</div>
		</div>
    )
  }
}

export default App;
