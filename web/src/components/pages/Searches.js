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

    fetch('https://wpdirectory.net/api/v1/searches/')
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
      <div className="page page-repos">
        <div className="panel plugins">
          <h2>Plugins Repository Overview</h2>
          <p>Below is a general overview of the data stored for WordPress plugins.</p>
          <ul className="details">
            <li><span className="name">Revision</span> {this.state.plugins.revision}</li>
            <li><span className="name">Total</span> {this.state.plugins.total}</li>
            <li><span className="name">Closed</span> {this.state.plugins.closed}</li>
            <li><span className="name">Pending Updates</span> {this.state.plugins.queue}</li>
          </ul>
        </div>
        <div className="panel themes">
          <h2>Themes Repository Overview</h2>
          <p>Below is a general overview of the data stored for WordPress themes.</p>
          <ul className="details">
            <li><span className="name">Revision</span> {this.state.themes.revision}</li>
            <li><span className="name">Total</span> {this.state.themes.total}</li>
            <li><span className="name">Closed</span> {this.state.themes.closed}</li>
            <li><span className="name">Pending Updates</span> {this.state.themes.queue}</li>
          </ul>
        </div>
      </div>
    )
  }
}

export default Searches