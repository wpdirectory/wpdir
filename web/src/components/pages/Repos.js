import React, { Component } from 'react'
import Loadicon from '../general/Loadicon.js'
import Hostname from '../../utils/Hostname.js'

class Repos extends Component {

  constructor(props) {
    super(props);
    this.state = {
        plugins: '',
        themes: '',
        isLoading: true,
    }
  }

  componentWillMount = () => {
    fetch( Hostname + '/api/v1/repos/overview' )
    .then( response => {
      return response.json()
    })
    .then( data => {
      this.setState({
        plugins: data.plugins,
        themes: data.themes,
        isLoading: false
      })
    })
  }

  componentDidMount() {
    document.title = 'Repositories Overview - WPdirectory'
  }

  render() {
    if (this.state.isLoading === true) {
      return (
        <div className="page page-repos grid-container">
          <div className="grid-x grid-margin-x grid-margin-y">
            <div className="panel cell small-12 medium-6">
              <h2>Plugins Overview</h2>
              <p>Below is a general overview of the data stored for WordPress plugins.</p>
              <Loadicon />
            </div>
            <div className="panel cell small-12 medium-6">
              <h2>Themes Overview</h2>
              <p>Below is a general overview of the data stored for WordPress themes.</p>
              <Loadicon />
            </div>
          </div>
        </div>
      )
    }

    return (
      <div className="page page-repos grid-container">
        <div className="grid-x grid-margin-x grid-margin-y">
          <div className="panel cell small-12 medium-6">
            <h2>Plugins Overview</h2>
            <p>Below is a general overview of the data stored for WordPress plugins.</p>
            <ul className="details">
              <li><span className="name">Revision</span> {this.state.plugins.revision}</li>
              <li><span className="name">Total</span> {this.state.plugins.total}</li>
              <li><span className="name">Pending Updates</span> {this.state.plugins.queue}</li>
            </ul>
          </div>
          <div className="panel cell small-12 medium-6">
            <h2>Themes Overview</h2>
            <p>Below is a general overview of the data stored for WordPress themes.</p>
            <ul className="details">
              <li><span className="name">Revision</span> {this.state.themes.revision}</li>
              <li><span className="name">Total</span> {this.state.themes.total}</li>
              <li><span className="name">Pending Updates</span> {this.state.themes.queue}</li>
            </ul>
          </div>
        </div>
      </div>
    )
  }
}

export default Repos