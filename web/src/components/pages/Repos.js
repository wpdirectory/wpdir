import React, { Component } from 'react'

class Repos extends Component {

  constructor(props) {
    super(props);
    this.state = {
        plugins: '',
        themes: '',
    };
  }

  componentWillMount = () => {

    fetch('https://wpdirectory.net/api/v1/repos/overview/')
    .then( response => {
      return response.json()
    })
    .then( data => {
      this.setState({plugins: data.plugins})
      this.setState({themes: data.themes})
    })

  }

  componentDidMount() {
    document.title = 'Repositories - WPdirectory'
  }

  render() {

    return (
      <div className="page page-repos">
        <div className="panel plugins">
          <h2>Plugins Repository Overview</h2>
          <p>Below is a general overview of the data stored for WordPress plugins.</p>
          <ul className="details">
            <li><span className="name">Revision</span> {this.state.plugins.revision}</li>
            <li><span className="name">Open</span> {this.state.plugins.total - this.state.plugins.closed}</li>
            <li><span className="name">Closed</span> {this.state.plugins.closed}</li>
            <li><span className="name">Total</span> {this.state.plugins.total}</li>
            <li><span className="name">Pending Updates</span> {this.state.plugins.queue}</li>
          </ul>
        </div>
        <div className="panel themes">
          <h2>Themes Repository Overview</h2>
          <p>Below is a general overview of the data stored for WordPress themes.</p>
          <ul className="details">
            <li><span className="name">Revision</span> {this.state.themes.revision}</li>
            <li><span className="name">Open</span> {this.state.themes.total - this.state.themes.closed}</li>
            <li><span className="name">Closed</span> {this.state.themes.closed}</li>
            <li><span className="name">Total</span> {this.state.themes.total}</li>
            <li><span className="name">Pending Updates</span> {this.state.themes.queue}</li>
          </ul>
        </div>
      </div>
    )
  }
}

export default Repos