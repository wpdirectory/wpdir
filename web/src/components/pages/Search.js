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

  tracURL = (slug, name, linenum) => {
    console.log(name)
    let len = slug.length
    name = name.slice((len * 2) + 1)
    return 'https://plugins.trac.wordpress.org/browser/' + slug + '/trunk' + name + '/#L' + linenum
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
          <div className="search-info panel">
            <h2>Overview</h2>
            <div className="info">
              <dl>
                <dt>Search Input</dt>
                <dd>{this.state.input}</dd>
                <dt>Repository</dt>
                <dd>{this.upperCaseFirst(this.state.repo)}</dd>
                <dt>Matches</dt>
                <dd>{this.state.matches}</dd>
                <dt>Time Taken</dt>
                <dd>{duration}</dd>
              </dl>
            </div>
          </div>
        )
      case 1:
        return (
          <div className="search-info panel">
            <h2>Overview</h2>
            <div className="info">
              <dl>
                <dt>Search Input</dt>
                <dd>{this.state.input}</dd>
                <dt>Repository</dt>
                <dd>{this.upperCaseFirst(this.state.repo)}</dd>
                <dt>Matches</dt>
                <dd>{this.state.matches}</dd>
                <dt>Progress</dt>
                <dd>{this.calcProgress(this.state.progress, this.state.total)}%</dd>
                <dt>Duration</dt>
                <dd>{duration}</dd>
              </dl>
            </div>
          </div>
        )
      default:
        return (
          <div className="search-info panel">
            <h2>Overview</h2>
            <div className="info">
              <dl>
                <dt>Search Input</dt>
                <dd>{this.state.input}</dd>
                <dt>Repository</dt>
                <dd>{this.upperCaseFirst(this.state.repo)}</dd>
              </dl>
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
        <div className="search-summary panel">
          <h2>Summary <small> {this.state.summary.total}{ ' matches'}</small></h2>
          <Summary repo={this.state.repo} id={this.state.id} items={this.state.summary.list} />
        </div>
      )
    } else {
      searchSummary = (
        <div className="search-summary panel">
          <Loadicon />
        </div>
      )
    }

    if ( this.state.isLoading === true ) {
      return (
        <div className="page page-search">
          <div className="title panel">
            <h1>Search</h1>
          </div>
          <div className="search-info panel">
            <Loadicon />
          </div>
        </div>
      )
    } else {
      return (
        <div className="page page-search">
          <div className="title panel">
            <h1 style={margin}>Search - {this.getStatus(this.state.status)}</h1>
            {(() => {
              if (this.state.status === 1) {
                return (<ProgressBar progress={this.calcProgress(this.state.progress, this.state.total)} />)
              }
            })()}
          </div>
          {searchSummary}
          {this.formatOverview()}
        </div>
      )
    }
  }
}

export default Search