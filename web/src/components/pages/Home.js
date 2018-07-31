import React, { Component } from 'react'
import { Link } from 'react-router-dom';
import SearchForm from '../general/search/SearchForm'
import Dashicon from '../general/Dashicon.js'

class Home extends Component {

  constructor(props) {
    super(props);
    this.state = {
        searches: [],
    };
  }

  componentWillMount = () => {

    fetch(window.wpdirHost + '/api/v1/searches/10')
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      if (data.searches) {
        this.setState({searches: data.searches})
      }
    })

  }

  getRepoIcon = (repo) => {
    switch(repo) {
      case 'themes':
          return <Dashicon icon="admin-appearance" size={ 22 } />
      default:
          return <Dashicon icon="admin-plugins" size={ 22 } />
    }
  }

  componentDidMount() {
    document.title = 'WordPress Directory Searcher - WPdirectory'
  }

  render() {
    console.log(this.state.searches)
    let latestSearches

    if ( this.state.searches.length && this.state.searches.length > 0 ) {
      let searches = this.state.searches.sort( (a,b) => Date.parse(a.started) < Date.parse(b.started) )
      latestSearches = searches.map( (search, idx) => {
        return (
          <li key={idx}>
            <span className="input"><Link to={'/search/' + search.id} title={search.input}>{search.input.substring(0, 34)}</Link></span>
            <span className="matches">{search.matches}</span>
            <span className="directory" title={search.repo.charAt(0).toUpperCase() + search.repo.slice(1)}>{this.getRepoIcon(search.repo)}</span>
          </li>
        )
      })

    } else {

      latestSearches = <p>Sorry, no searches found.</p>

    }

    return (
      <div className="page page-home grid-container">
        <div className="grid-x grid-margin-x grid-margin-y">
          <div className="panel cell small-12 medium-6 large-4">
            <SearchForm />
          </div>
          <div className="panel cell small-12 medium-6 large-4">
            <h3>Search Tips</h3>
            <p>The search input uses <a href="https://github.com/google/re2/wiki/Syntax" target="_blank" rel="noopener noreferrer">RE2</a> regex and may use syntax a little different to what you are used to.</p>
            <p>Here are a few examples to help get you started:</p>
            <pre>
              Searching for a function<br />
              <code>register_meta\(</code>

              <br /><br />TODO: Add more examples.
            </pre>
          </div>
          <div className="panel cell small-12 medium-12 large-4">
            <h3>Recent Searches</h3>
            <ul className="search-list">
              {latestSearches}
            </ul>
          </div>
        </div>
      </div>
    )
  }
}

export default Home