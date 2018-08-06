import React, { Component } from 'react'
import Match from './Match.js'
import API from '../../../utils/API.js'

class Matches extends Component {

  constructor(props) {
    super(props)
    this.state = {
      matches: [],
      isLoading: true,
      error: ''
    }
  }

  componentWillMount = () => {
    this.setState({ isLoading: true })

    API.get( '/search/matches/' + this.props.id + '/' + this.props.slug )
      .then( result => this.setState({
        matches: result.data.list,
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
  }

  render() {
    const { 
      matches,
      isLoading,
      error
    } = this.state

    let matchList

    if ( isLoading ) {
      matchList = <p>Loading matches...</p>
    }  else {
      if ( error ) {
        matchList = <p className="error">Sorry, there was a problem fetching data.</p>
      } else {
        matchList = matches.map( (match, key) => {
          return (
            <Match repo={this.props.repo} match={match} key={key} />
          )
        })
      }
    }

    return (
      <ul className="matches">
        {matchList}
      </ul>
	  )
  }
}

export default Matches