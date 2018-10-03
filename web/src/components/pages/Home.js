import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import SearchForm from '../general/search/SearchForm'
import Dashicon from '../general/Dashicon.js'
import Loadicon from '../general/Loadicon.js'
import API from '../../utils/API.js'

class Home extends Component {

  constructor(props) {
    super(props)
    this.state = {
        searches: [],
        isLoading: true,
        error: '',
    }
  }

  componentWillMount = () => {
    this.setState({ isLoading: true })

    API.get( '/searches/10' )
      .then( result => this.setState({
        searches: result.data.searches,
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
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
    const { 
      searches,
      isLoading,
      error
    } = this.state

    let latestSearches

    if ( isLoading ) {
      latestSearches = <Loadicon />
    }  else {
      if ( !Array.isArray(searches) || !searches.length || error ) {
        latestSearches = <p className="error">Sorry, there was a problem fetching data.</p>
      } else {
        latestSearches = searches.map( (search, idx) => {
          return (
            <li key={idx}>
              <span className="input"><Link to={'/search/' + search.id} title={search.input}>{search.input.substring(0, 34)}</Link></span>
              <span className="matches">{search.matches}</span>
              <span className="directory" title={search.repo.charAt(0).toUpperCase() + search.repo.slice(1)}>{this.getRepoIcon(search.repo)}</span>
            </li>
          )
        })
      }
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
            <p>Using regular expressions (regex) can be difficult so the <Link to={'/examples/'} title={'Examples'}>examples</Link> page guides you through some common searches and explains how they match the intended targets.</p>
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