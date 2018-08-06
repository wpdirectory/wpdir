import React, { Component } from 'react'
import Loadicon from '../Loadicon.js'
import API from '../../../utils/API.js'

class Overview extends Component {

  constructor(props) {
    super(props)
    this.state = {
      interval: 0,
      id: this.props.id,
      input: '',
      repo: '',
      started: 0,
      completed: 0,
      progress: 0,
      status: 5,
      matches: 0,
      isLoading: true,
      error: '',
    }
  }

  componentWillMount = () => {
    this.fetchData()
  }

  fetchData = () => {
    this.setState({ isLoading: true })

    API.get( '/search/' + this.state.id )
      .then( result => this.setState({
        id: result.data.id,
        input: result.data.input,
        repo: result.data.repo,
        progress: result.data.progress,
        status: result.data.status,
        matches: result.data.matches,
        started: Date.parse(result.data.started),
        completed: ( result.data.completed ? Date.parse(result.data.completed) : 0 ),
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
  }

  render() {
    const { 
      isLoading,
      error
    } = this.state

    if ( isLoading ) {
      latestSearches = <Loadicon />
    }  else {
      if ( error ) {
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
      <div>
        
      </div>
	  )
  }
}

export default Overview