import React, { Component } from 'react'
import Dashicon from '../general/Dashicon.js'

class Search extends Component {

  constructor(props) {
    super(props);
    this.state = {
        id: '',
        results: []
    };
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/search/' + this.props.match.params.id + '/')
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({id: data.id})
      if (data.results) {
        this.setState({results: data.results})
      }
    })

  }

  componentDidMount(){
    document.title = 'Search: ' + this.state.id + ' - WPdirectory'
  }

  render() {

    let searchResults

    if ( this.state.results.length && this.state.results.length > 0 ) {

    searchResults = this.state.results.map( function(result, idx) {
      let slug = result.Slug
      return (
        <div className="" key={idx}>
      {result.Matches.map(function(match, idx) {
        let filename = match.Filename
        return (
          <div className="" key={idx}>
          {console.log(match)}
          {match.Matches.map(function(m, idx) {
            return (
            <div key={idx} className="result">
            <div className="file">
              <span className="name">{filename}</span>
              <a className="link" href={"https://plugins.trac.wordpress.org/browser/" + slug} target="_blank" rel="noopener noreferrer">
                <Dashicon icon="external" size={ 22 } />
              </a>
            </div>
            <ul className="lines">
                <li>
                  <span className="num">{m.LineNumber}</span>
                  <span className="excerpt"><code>{m.Line}</code></span>
                </li>
            </ul>
          </div>
            )
          })}
          </div>
        )
      })}
      </div>
      )
    })

  } else {

    searchResults = <p>Sorry, no results found.</p>

  }

      return (
        <div className="page page-search">
          <div className="title panel">
            <h1>Search</h1>
          </div>
          <div className="search-info panel">
            <h2>Overview</h2>
            <div className="info">
              <dl>
                <dt>Total Matches</dt>
                <dd>317</dd>
              </dl>
            </div>
          </div>
          <div className="search-results panel">
            <h2>Results</h2>
            <div className="results">
              {searchResults}
            </div>
          </div>
        </div>
      )
  }
}

export default Search