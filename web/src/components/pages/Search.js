import React, { Component } from 'react'
import Loadicon from '../general/Loadicon.js'
import ProgressBar from '../general/ProgressBar.js'
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
    };
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
      this.setState({ total: data.total })
      this.setState({ status: data.status })
      this.setState({ matches: data.matches })
      if (data.summary) {
        this.setState({ summary: data.summary })
      }
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
      this.setState({ total: data.total })
      this.setState({ status: data.status })
      this.setState({ matches: data.matches })
      if (data.summary) {
        this.setState({ summary: data.summary })
      }
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

  calcProgress = (current, total) => {
    if (current === 0) {
      return 0
    } else {
      return (( current / total ) * 100).toFixed(0)
    }
  }

  componentDidMount = () => {
    document.title = 'Search ' + this.state.id + ' - WPdirectory'
    this.interval = setInterval(() => {
      if ( this.state.status === 1 ) {
        this.refreshData()
      }
    }, 2000);
  }

  componentWillUnmount = () => {
    clearInterval(this.interval);
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
                <div className="cell small-12 medium-6">
                  <h5>Search Regex</h5>
                  {this.state.input}
                </div>
                <div className="cell small-12 medium-6">
                  <h5>Repository</h5>
                  {this.upperCaseFirst(this.state.repo)}
                </div>
                <div className="cell small-12 medium-6">
                  <h5>Total Matches</h5>
                  {this.state.matches}
                </div>
                <div className="cell small-12 medium-6">
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
              <div className="cell small-12 medium-6">
                <h5>Search Regex</h5>
                {this.state.input}
              </div>
              <div className="cell small-12 medium-6">
                <h5>Repository</h5>
                {this.upperCaseFirst(this.state.repo)}
              </div>
              <div className="cell small-12 medium-6">
                <h5>Total Matches</h5>
                {this.state.matches}
              </div>
              <div className="cell small-12 medium-6">
                <h5>Progress</h5>
                {this.calcProgress(this.state.progress, this.state.total)}%
              </div>
              <div className="cell small-12 medium-6">
                <h5>Duration</h5>
                {duration}
              </div>
            </div>
          </div>
        )
      default:
        return (
          <div className="search-info panel cell small-12">
            <h2>Info</h2>
            <div className="info grid-x grid-margin-x grid-margin-y">
              <div className="cell small-12 medium-6">
                <h5>Search Input</h5>
                {this.state.input}
              </div>
              <div className="cell small-12 medium-6">
                <h5>Repository</h5>
                {this.upperCaseFirst(this.state.repo)}
              </div>
            </div>
          </div>
        )
    }
  }

  render() {
    let margin
    if (this.state.status === 2) {
      margin = {
        margin: '0 0 0 0',
      }
    } else {
      margin = {}
    }

    let searchSummary
    if ( !!this.state.summary.list && this.state.summary.list.length && this.state.summary.list.length > 0 ) {
      searchSummary = (
        <div className="search-summary panel cell small-12">
          <h2>Summary <small>({this.state.summary.total}{ ' matches'})</small></h2>
          <Summary repo={this.state.repo} id={this.state.id} items={this.state.summary.list} />
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
          <div className="title panel cell small-12">
            <h1>Search</h1>
          </div>
          <div className="search-info panel cell small-12">
            <Loadicon />
          </div>
        </div>
      )
    } else {
      return (
        <div className="page page-search grid-container">
          <div className="grid-x grid-margin-x grid-margin-y">
            <div className="title panel cell small-12">
              <h1 style={margin}>Search - {this.getStatus(this.state.status)}</h1>
              {(() => {
                if (this.state.status === 1) {
                  return (<ProgressBar progress={this.calcProgress(this.state.progress, this.state.total)} />)
                }
              })()}
            </div>
            {this.formatOverview()}
            {searchSummary}
          </div>
        </div>
      )
    }
  }
}

export default Search