import React, { Component } from 'react'
import Match from './Match.js'
import Hostname from '../../../utils/Hostname.js'

class Matches extends Component {

  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      matches: [],
    }
  }

  componentWillMount = () => {
    fetch( Hostname + '/api/v1/search/matches/' + this.props.id + '/' + this.props.slug )
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({ matches: data.list })
      this.setState({ isLoading: false })
    })
  }

  render() {
    let matchList
    if ( !!this.state.matches && this.state.matches.length && this.state.matches.length > 0 ) {
      matchList = this.state.matches.map( (match, key) => {
        return (
          <Match repo={this.props.repo} match={match} key={key} />
        );
      })
    } else {
      matchList = <p>Loading matches...</p>
    }

    return (
      <ul className="matches">
        {matchList}
      </ul>
	  );
  }
}

export default Matches