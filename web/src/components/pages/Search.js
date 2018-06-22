import React, { Component } from 'react'
import Dashicon from '../general/Dashicon.js'

class Search extends Component {

  constructor(props) {
    super(props);
    this.state = {
      interval: 0,
      id: '',
      input: '',
      repo: '',
      started: '',
      completed: '',
      progress: 0,
      total: 0,
      status: 0,
      matches: []
    };
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/search/' + this.props.match.params.id + '/')
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({id: data.id})
      this.setState({input: data.input})
      this.setState({repo: data.repo})
      this.setState({started: Date.parse(data.started)})
      this.setState({completed: Date.parse(data.completed)})
      this.setState({progress: data.progress})
      this.setState({total: data.total})
      this.setState({status: data.status})
      if (data.matches) {
        this.setState({matches: data.matches})
      }
    })

  }

  tick = () => {
    fetch('https://wpdirectory.net/api/v1/search/' + this.props.match.params.id + '/')
    .then( response => {
      return response.json()
      
    })
    .then( data => {
      this.setState({id: data.id})
      this.setState({input: data.input})
      this.setState({repo: data.repo})
      this.setState({started: Date.parse(data.started)})
      this.setState({completed: Date.parse(data.completed)})
      this.setState({progress: data.progress})
      this.setState({total: data.total})
      this.setState({status: data.status})
      if (data.matches) {
        this.setState({matches: data.matches})
      }
    })
  }

  getStatus = (code) => {
    switch (code) {
      case 1:
        return 'Started'
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

  formatProgress = (current, total) => {
    return (( current / total ) * 100).toFixed(0) + '%'
  }

  componentDidMount() {
    document.title = 'Search ' + this.state.id + ' - WPdirectory'
    this.interval = setInterval(() => {
      if ( this.state.status !== 2 ) {
        this.tick()
      }
    }, 3000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  render() {

    let searchResults

    if ( this.state.matches.length && this.state.matches.length > 0 ) {

      searchResults = this.state.matches.map( (match, idx) => {
        return (
          <div key={idx} className="result">
            <div className="file">
              <span className="name" title={match.file}>{this.formatFilename(match.slug, match.file)}</span>
              <a className="link" href={this.tracURL(match.slug, match.file, match.line_num)} target="_blank" rel="noopener noreferrer">
                <Dashicon icon="external" size={ 22 } />
              </a>
            </div>
            <ul className="lines">
              <li><span className="num">{match.line_num - 2}</span><span className="excerpt"><code>{match.before[0]}</code></span></li>
              <li><span className="num">{match.line_num - 1}</span><span className="excerpt"><code>{match.before[1]}</code></span></li>
              <li><span className="num">{match.line_num}</span><span className="excerpt"><code>{match.line_text}</code></span></li>
              <li><span className="num">{match.line_num + 1}</span><span className="excerpt"><code>{match.after[0]}</code></span></li>
              <li><span className="num">{match.line_num + 2}</span><span className="excerpt"><code>{match.after[1]}</code></span></li>
            </ul>
          </div>
          )
      })

    } else {

      searchResults = <p>Sorry, no results found.</p>

    }

    let duration

    if (this.state.status !== 0) {
      duration = this.progressTime(this.state.started);
    } else {
      duration = 'In Queue'
    }

    return (
      <div className="page page-search">
        <div className="title panel">
          <h1>Search - {this.getStatus(this.state.status)} - {this.formatProgress(this.state.progress, this.state.total)}</h1>
        </div>
        <div className="search-info panel">
          <h2>Overview</h2>
          <div className="info">
            <dl>
              <dt>Search Input</dt>
              <dd>{this.state.input}</dd>
              <dt>Repository</dt>
              <dd>{this.upperCaseFirst(this.state.repo)}</dd>
              <dt>Matches</dt>
              <dd>{this.state.matches.length}</dd>
              <dt>Duration</dt>
              <dd>{duration}</dd>
            </dl>
           </div>
        </div>
        <div className="search-results panel">
          <h2>Results</h2>
          <div className="results">
            {searchResults}
          </div>
        </div>
      </div>
    )
  }
}

export default Search