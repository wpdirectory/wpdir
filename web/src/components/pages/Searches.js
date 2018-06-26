import React, { Component } from 'react'
import { Link } from 'react-router-dom';

class Searches extends Component {

  constructor(props) {
    super(props);
    this.state = {
      isLoading: true,
      searches: [],
    };
  }

  upperCaseFirst = (name) => {
    return name.charAt(0).toUpperCase() + name.slice(1)
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/searches/100')
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
    if ( !!this.state.searches && this.state.searches.length && this.state.searches.length > 0 ) {
      let searches = this.state.searches.sort( (a,b) => Date.parse(a.started) < Date.parse(b.started) );
      searchList = searches.map( (search, idx) => {
        return (
          <tr key={idx}><td><Link to={'/search/' + search.id}>{search.input}</Link></td><td>{this.upperCaseFirst(search.repo)}</td><td>{search.matches}</td></tr>
        )
      })
    } else {
      searchList = <tr><th>Sorry, no searches found.</th></tr>
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