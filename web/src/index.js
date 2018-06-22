import React from 'react';
import ReactDOM from 'react-dom';

import './assets/scss/app.scss'

import App from './App';
//import registerServiceWorker from './registerServiceWorker';
import { unregister } from './registerServiceWorker';

ReactDOM.render(<App />, document.getElementById('root'));

//registerServiceWorker();
unregister();