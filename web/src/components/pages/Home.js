import React, { Component } from 'react'
import SearchForm from '../general/SearchForm'

class Home extends Component {

  componentDidMount(){
    document.title = 'WordPress Directory Searcher - WPdirectory'
  }

  render() {
      return (
        <div className="page page-home">
          <div className="grid-x gutter-x gutter-y small-1 medium-2 large-3">
            <div className="panel">
              <SearchForm />
            </div>
            <div className="panel">
              <h3>Recent Searches</h3>
              <ul className="search-list">
                <li>
                  <span className="input">regex input</span>
                  <span className="matches">712</span>
                  <span className="directory">Plugins</span>
                </li>
                <li>
                  <span className="input">regex input</span>
                  <span className="matches">712</span>
                  <span className="directory">Plugins</span>
                </li>
                <li>
                  <span className="input">regex input</span>
                  <span className="matches">712</span>
                  <span className="directory">Plugins</span>
                </li>
              </ul>
            </div>
            <div className="panel">
              <h3>Overview</h3>
              <dl>
                <dt>Total Plugins</dt>
                <dd>74,746</dd>
                <dt>Total Themes</dt>
                <dd>13,945</dd>
              </dl>
            </div>
          </div>
        </div>
      )
  }
}

export default Home