import React, { Component } from 'react'
import Loadicon from '../general/Loadicon.js'
import ProgressBlock from '../general/ProgressBlock.js'
import Summary from '../general/search/Summary.js'

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: true,
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
      summary: {
        list: [],
        total: 0,
      },
    }
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/search/' + this.props.match.params.id)
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({ id: data.id })
      this.setState({ input: data.input })
      this.setState({ repo: data.repo })
      if (data.started) {
        this.setState({ started: Date.parse(data.started) })
      }
      if (data.completed) {
        this.setState({ completed: Date.parse(data.completed) })
      }
      this.setState({ progress: data.progress })
      this.setState({ status: data.status })
      this.setState({ matches: data.matches })
      this.setState({ isLoading: false })
    })

  }

  refreshData = () => {
    fetch('https://wpdirectory.net/api/v1/search/' + this.props.match.params.id)
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({ id: data.id })
      this.setState({ input: data.input })
      this.setState({ repo: data.repo })
      if (data.started) {
        this.setState({ started: Date.parse(data.started) })
      }
      if (data.completed) {
        this.setState({ completed: Date.parse(data.completed) })
      }
      this.setState({ progress: data.progress })
      this.setState({ status: data.status })
      this.setState({ matches: data.matches })
      this.setState({ isLoading: false })
    })
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

  progressTime = (started) => {
    let active
    if (this.state.completed > 0) {
      active = this.state.completed - started
    } else {
      active = Date.now() - started
    }
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
        this.refreshData()
      }
    }, 5000)
    this.updateInterval = setInterval(() => {
      if ( this.state.status === 1 ) {
        this.refreshData()
      }
    }, 2000)
  }

  componentWillUnmount = () => {
    clearInterval(this.updateInterval)
    clearInterval(this.queueInterval)
  }

  formatOverview = () => {
    let duration
    if (this.state.started > 0) {
      duration = this.progressTime(this.state.started);
    }
    switch( this.state.status ) {
      case 2:
        return (
          <div className="search-info panel cell small-12">
            <h2>Overview</h2>
            <div className="info">
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
                    if (this.state.matches > 10000) {
                      return (<label className="is-invalid-label">Search aborted after hitting match limit (10,000).</label>)
                    }
                  })()}
                </div>
                <div className="cell small-12 medium-4">
                  <h5>Time Taken</h5>
                  {duration}
                </div>
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
    let searchSummary
    if ( this.state.status === 2 ) {
      searchSummary = (
        <div className="search-summary panel cell small-12">
          <h2>Summary <small>({this.state.matches}{ ' matches'})</small></h2>
          <Summary repo={this.state.repo} id={this.state.id} matches={this.state.matches} />
        </div>
      )
    } else {
      searchSummary = (
        <div className="search-summary panel cell small-12">
          <Loadicon />
        </div>
      )
    }

    if ( this.state.isLoading === true ) {
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
                return searchSummary
              }
            })()}
          </div>
        </div>
      )
    }
  }
}

export default Search