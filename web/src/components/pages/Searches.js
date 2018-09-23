import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import Loadicon from '../general/Loadicon.js'
import API from '../../utils/API.js'

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
    this.setState({ isLoading: true })

    API.get( '/searches/100' )
      .then( result => this.setState({
        searches: result.data.searches,
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
  }

  componentDidMount() {
    document.title = 'Searches - WPdirectory'
  }

  upperCaseFirst = (name) => {
    return name.charAt(0).toUpperCase() + name.slice(1)
  }

  render() {
    const { 
      searches,
      isLoading,
      error
    } = this.state

    let searchList

    if ( isLoading ) {
      searchList = <tr><td><Loadicon /></td></tr>
    }  else {
      if ( !Array.isArray(searches) || !searches.length || error ) {
        searchList = <tr><td><p className="error">Sorry, there was a problem fetching data.</p></td></tr>
      } else {
        searchList = searches.map( (search, idx) => {
          return (
            <tr key={idx}><td><Link to={'/search/' + search.id}>{search.input}</Link></td><td>{this.upperCaseFirst(search.repo)}</td><td>{search.matches}</td></tr>
          )
        })
      }
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