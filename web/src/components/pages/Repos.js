import React, { Component } from 'react'
import timeago from 'timeago.js'
import Loadicon from '../general/Loadicon.js'
import API from '../../utils/API.js'


class Repos extends Component {

  constructor(props) {
    super(props);
    this.state = {
        plugins: '',
        themes: '',
        pluginChartInstalls: '',
        themeChartInstalls: '',
        pluginChartSize: '',
        themeChartSize: '',
        isLoading: true,
        error: ''
    }
  }

  componentWillMount = () => {
    this.setState({ isLoading: true })

    API.get( '/repos/overview' )
      .then( result => this.setState({
        plugins: result.data.plugins,
        themes: result.data.themes,
        pluginChartInstalls: result.data.plugin_chart_installs,
        themeChartInstalls: result.data.theme_chart_installs,
        pluginChartSize: result.data.plugin_chart_size,
        themeChartSize: result.data.theme_chart_size,
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
  }

  componentDidMount() {
    document.title = 'Repositories Overview - WPdirectory'
  }

  render() {
    const {
      plugins,
      themes,
      pluginChartInstalls,
      themeChartInstalls,
      pluginChartSize,
      themeChartSize,
      isLoading,
      error
    } = this.state

    let pluginsContent
    let themesContent

    if ( isLoading ) {
      pluginsContent = <Loadicon />
      themesContent = <Loadicon />
    } else {
      if ( error ) {
        pluginsContent = <p className="error">Sorry, there was a problem fetching data.</p>
        themesContent = <p className="error">Sorry, there was a problem fetching data.</p>
      } else {
        pluginsContent = (
          <ul className="details">
            <li><span className="name">Revision</span> {plugins.revision}</li>
            <li><span className="name">Total</span> {plugins.total}</li>
            <li><span className="name">Open</span> {plugins.total - plugins.closed}</li>
            <li><span className="name">Closed</span> {plugins.closed}</li>
            <li><span className="name">Updated</span> <time dateTime={plugins.updated} title={plugins.updated}>{timeago().format(Date.parse(plugins.updated))}</time></li>
          </ul>
        )
        themesContent = (
          <ul className="details">
            <li><span className="name">Revision</span> {themes.revision}</li>
            <li><span className="name">Total</span> {themes.total}</li>
            <li><span className="name">Open</span> {themes.total - themes.closed}</li>
            <li><span className="name">Closed</span> {themes.closed}</li>
            <li><span className="name">Updated</span> <time dateTime={themes.updated} title={themes.updated}>{timeago().format(Date.parse(themes.updated))}</time></li>
          </ul>
        )
      }
    }

    return (
      <div className="page page-repos grid-container">
        <div className="grid-x grid-margin-x grid-margin-y">
          <div className="panel cell small-12 medium-12 large-6">
            <h2>Plugins Overview</h2>
            <p>Below is a general overview of the data stored for WordPress plugins.</p>
            <div className="stats">
              {pluginsContent}
            </div>
            <div className="chart-installs">
              <h3>Installs</h3>
              <div dangerouslySetInnerHTML={{__html: pluginChartInstalls}} />
            </div>
            <div className="chart-size">
              <h3>Files / Size (MB)</h3>
              <div dangerouslySetInnerHTML={{__html: pluginChartSize}} />
            </div>
          </div>
          <div className="panel cell small-12 medium-12 large-6">
            <h2>Themes Overview</h2>
            <p>Below is a general overview of the data stored for WordPress themes.</p>
            <div className="stats">
              {themesContent}
            </div>
            <div className="chart-installs">
              <h3>Installs</h3>
              <div dangerouslySetInnerHTML={{__html: themeChartInstalls}} />
            </div>
            <div className="chart-size">
              <h3>Files / Size (MB)</h3>
              <div dangerouslySetInnerHTML={{__html: themeChartSize}} />
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default Repos