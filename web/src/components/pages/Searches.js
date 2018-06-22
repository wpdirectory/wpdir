import React, { Component } from 'react'
import Dashicon from '../general/Dashicon.js'

class Searches extends Component {

  constructor(props) {
    super(props);
    this.state = {
        searches: []
    };
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/searches')
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({searches: data.searches})
    })

  }

  componentDidMount() {
    document.title = 'Searches - WPdirectory'
  }

  render() {

    let searchList
    if ( this.state.searches.length && this.state.searches.length > 0 ) {
      searchList = this.state.matches.map( function(match, idx) {
        return (
          <div key={idx} className="result">
            <div className="file">
              <span className="name">{match.file}</span>
              <a className="link" href={"https://plugins.trac.wordpress.org/browser/" + match.slug} target="_blank" rel="noopener noreferrer">
                <Dashicon icon="external" size={ 22 } />
              </a>
            </div>
            <ul className="lines">
                <li>
                  <span className="num">{match.line_num}</span>
                  <span className="excerpt"><code>{match.line_text}</code></span>
                </li>
            </ul>
          </div>
        )
      })
    } else {
      searchList = <p>Sorry, no searches found.</p>
    }

    return (
      <div className="page page-searches">
        <div className="panel searches">
          <h2>Search List</h2>
          <ul className="details">
            {searchList}
          </ul>
        </div>
      </div>
    )
  }
}

export default Searches