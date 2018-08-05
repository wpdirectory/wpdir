import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import Loadicon from '../general/Loadicon.js'
import Hostname from '../../utils/Hostname.js'

class Searches extends Component {

  constructor(props) {
    super(props)
    this.state = {
      searches: [],
      isLoading: true,
      error: '',
    }
  }

  componentWillMount = () => {
    fetch( Hostname + '/api/v1/searches/100' )
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({
        searches: data.searches,
        isLoading: false
      })
    })
    .catch( error => {
      this.setState({
        isLoading: false,
        error: error
      })
    })
  }

  componentDidMount() {
    document.title = 'Searches - WPdirectory'
  }

  upperCaseFirst = (name) => {
    return name.charAt(0).toUpperCase() + name.slice(1)
  }

  render() {
    let searchList
    if ( !!this.state.searches && this.state.searches.length && this.state.searches.length > 0 ) {
      searchList = this.state.searches.map( (search, idx) => {
        return (
          <tr key={idx}><td><Link to={'/search/' + search.id}>{search.input}</Link></td><td>{this.upperCaseFirst(search.repo)}</td><td>{search.matches}</td></tr>
        )
      })
    } else {
      searchList = <tr><th>Sorry, no searches found.</th></tr>
    }

    if (this.state.isLoading === true) {
      searchList = <tr><td><Loadicon /></td></tr>
    }
    
    return (
      <div className="page page-searches grid-container">
        <div className="grid-x grid-margin-x grid-margin-y">
          <div className="panel cell small-12">
            <h2>Search List</h2>
            <table className="searches-table">
              <thead>
                <tr>
                  <th width="auto">Input</th>
                  <th width="100">Repo</th>
                  <th width="100">Matches</th>
                </tr>
              </thead>
              <tbody>
                {searchList}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    )
  }
}

export default Searches