import React, { Component } from 'react'
import { BrowserRouter as Router, Switch, Route, NavLink, Link } from 'react-router-dom'

import './assets/scss/app.scss'

//import Header from './components/Header'
import Footer from './components/Footer'
import Home from './components/pages/Home'
import NotFound from './components/pages/NotFound'

class App extends Component {

  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Router>
        <div className="app">
        
          <header className="header grid-x gutter-x">
            <div className="title">
              <Link to="/">WPDir</Link>
            </div>
            <nav className="main-menu" aria-label="Main Navigation">
              <ul className="menu">
                <li><NavLink to="/searches">Searches</NavLink></li>
                <li><NavLink to="/search/new">New Search</NavLink></li>
                <li><NavLink to="/about">About</NavLink></li>
              </ul>
            </nav>
          </header>

          <section className="content grid-x padding-y">
            <Switch>
              <Route exact path="/" component={Home} />
              <Route path="/searches" component={Home} />
              <Route path="/search/new" component={Home} />
              <Route path="/search/:id" component={Home} />
              <Route path="/about" component={Home} />
              <Route component={NotFound} />
            </Switch>
          </section>

          <Footer />

        </div>
      </Router>
    )
  }

}

export default App;