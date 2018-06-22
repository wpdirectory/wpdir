import React, { Component } from 'react'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'

import './assets/scss/app.scss'

import Header from './components/Header'
import Footer from './components/Footer'
import Home from './components/pages/Home'
import Search from './components/pages/Search'
import Repos from './components/pages/Repos'
import About from './components/pages/About'
import NotFound from './components/pages/NotFound'

class App extends Component {
  render() {
    return (
      <Router>
        <div className="app">
          <Header />

          <section className="content grid-x padding-y">
            <Switch>
              <Route exact path="/" component={Home} />
              <Route path="/searches" component={Home} />
              <Route path="/search/new" component={Home} />
              <Route path="/search/:id" component={Search} />
              <Route path="/repos" component={Repos} />
              <Route path="/about" component={About} />
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