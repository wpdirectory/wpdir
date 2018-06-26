import React, { Component } from 'react'
import Match from './Match.js'

class Matches extends Component {

  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      matches: [],
    }
  }

  componentWillMount = () => {
    fetch('https://wpdirectory.net/api/v1/search/matches/' + this.props.id + '/' + this.props.slug)
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      if (data.matches) {
        this.setState({ matches: data.matches })
      }
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