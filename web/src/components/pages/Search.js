import React, { Component } from 'react'
import Loadicon from '../general/Loadicon.js'
import ProgressBlock from '../general/ProgressBlock.js'
import Summary from '../general/search/Summary.js'
import API from '../../utils/API.js'

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      interval: 0,
      id: '',
      input: '',
      repo: '',
      started: 0,
      completed: 0,
      progress: 0,
      total: 0,
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

    API.get( '/search/' + this.props.match.params.id )
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

  getStatus = (code) => {
    switch (code) {
      case 1:
        return 'In Progress'
      case 2:
        return 'Completed'
      default:
        return 'Queued'
    }
  }

  upperCaseFirst = (name) => {
    return name.charAt(0).toUpperCase() + name.slice(1)
  }

  timeTaken = () => {
    let active = this.state.completed - this.state.started
    return Math.floor(active/1000) + ' Seconds'
  }

  formatFilename = (slug, name) => {
    let len = slug.length
    name = name.slice((len * 2) + 1)
    if (name.length > 30) {
      name = '...' + name.slice((38 - slug.length) - name.length)
    }

    return slug + '/' + name
  }

  formatName = (slug, name, version) => {
    if (name === undefined) {
      return slug
    } else {
      if (version === undefined) {
        return name
      } else {
        return name + ' (' + version + ')'
      }
    }
  }

  componentDidMount = () => {
    document.title = 'Search ' + this.state.id + ' - WPdirectory'
    this.queueInterval = setInterval(() => {
      if ( this.state.status === 0 ) {
        this.fetchData()
      }
    }, 5000)
    this.updateInterval = setInterval(() => {
      if ( this.state.status === 1 ) {
        this.fetchData()
      }
    }, 2000)
  }

  componentWillUnmount = () => {
    clearInterval(this.updateInterval)
    clearInterval(this.queueInterval)
  }

  formatOverview = () => {
    let duration
    if (this.state.status === 2) {
      duration = this.timeTaken()
    }
    switch( this.state.status ) {
      case 2:
        return (
          <div className="search-info panel cell small-12">
            <h2>Overview</h2>
              <div className="info grid-x grid-margin-x grid-margin-y">
                <div className="cell small-12">
                  <h5>Search Regex</h5>
                  <pre>{this.state.input}</pre>
                </div>
                <div className="cell small-12 medium-4">
                  <h5>Repository</h5>
                  {this.upperCaseFirst(this.state.repo)}
                </div>
                <div className="cell small-12 medium-4">
                  <h5>Total Matches</h5>
                  {this.state.matches}
                  {(() => {
                    if (this.state.matches > 100000) {
                      return (<label className="is-invalid-label">Search aborted after hitting match limit (100,000).</label>)
                    }
                  })()}
                </div>
                <div className="cell small-12 medium-4">
                  <h5>Time Taken</h5>
                  {duration}
                </div>
              </div>
          </div>
        )
      case 1:
        return (
          <div className="search-info panel cell small-12">
            <h2>Info</h2>
            <div className="info grid-x grid-margin-x grid-margin-y">
              <div className="cell small-12">
                <h5>Search Regex</h5>
                <pre>{this.state.input}</pre>
              </div>
              <div className="cell small-12 medium-4">
                <h5>Repository</h5>
                {this.upperCaseFirst(this.state.repo)}
              </div>
              <div className="cell small-12 medium-4">
                <h5>Total Matches</h5>
                {this.state.matches}
              </div>
            </div>
          </div>
        )
      default:
        return (
          <div className="search-info panel cell small-12">
            <h2>Info</h2>
            <div className="info grid-x grid-margin-x grid-margin-y">
              <div className="cell small-12">
                <h5>Search Input</h5>
                <pre>{this.state.input}</pre>
              </div>
              <div className="cell small-12 medium-4">
                <h5>Repository</h5>
                {this.upperCaseFirst(this.state.repo)}
              </div>
            </div>
          </div>
        )
    }
  }

  render() {
    const {
      isLoading,
      error
    } = this.state

    let summary
    if ( this.state.status === 2 ) {
      summary = (
        <div className="search-summary panel cell small-12">
          <h2>Summary <small>({this.state.matches}{ ' matches'})</small></h2>
          <Summary repo={this.state.repo} id={this.state.id} matches={this.state.matches} />
        </div>
      )
    } else {
      summary = (
        <div className="search-summary panel cell small-12">
          <Loadicon />
        </div>
      )
    }

    if ( isLoading ) {
      return (
        <div className="page page-search grid-container">
          <div className="grid-x grid-margin-x grid-margin-y">
            <div className="title panel cell small-12">
              <h1>Search</h1>
            </div>
            <div className="search-info panel cell small-12">
              <Loadicon />
            </div>
          </div>
        </div>
      )
    }  else {
      if ( error ) {
        return (
          <div className="page page-search grid-container">
            <div className="grid-x grid-margin-x grid-margin-y">
              <div className="title panel cell small-12">
                <h1>Search - {this.getStatus(this.state.status)}</h1>
              </div>
              <div className="search-info panel cell small-12">
                <p className="error">Sorry, there was a problem fetching data.</p>
              </div>
            </div>
          </div>
        )
      } else {
        return (
          <div className="page page-search grid-container">
            <div className="grid-x grid-margin-x grid-margin-y">
              <div className="title panel cell small-12">
                <h1>Search - {this.getStatus(this.state.status)}</h1>
              </div>
              {(() => {
                if (this.state.status === 1) {
                  return (<ProgressBlock progress={this.state.progress} status={this.getStatus(this.state.status)} />)
                }
              })()}
              {this.formatOverview()}
              {(() => {
                if (this.state.status === 2 && this.state.matches > 0) {
                  return summary
                }
              })()}
            </div>
          </div>
        )
      }
    }
  }
}

export default Search