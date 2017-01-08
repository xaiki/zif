require("babel-core/register");

import React from 'react';
import { render } from 'react-dom';
import injectTapEventPlugin from 'react-tap-event-plugin';

injectTapEventPlugin();

import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

import App from './components/App';

render((
	<MuiThemeProvider>
		<App/>
	</MuiThemeProvider>

), document.getElementById('root'))
