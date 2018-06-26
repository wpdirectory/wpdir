import React, { Component } from 'react'
import Highlighter from '../Highlighter.js'

class Modal extends Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      code: '',
      lang: '',
    }
  }

  getInitialState = () => {
		return {
			code: "// Code",
		}
  }

  componentWillMount = () => {
    fetch('https://wpdirectory.net/api/v1/file', {
      body: JSON.stringify({repo: this.props.repo, slug: this.props.match.slug, file: this.props.match.file}),
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'user-agent': 'WPDirectory/0.1.0',
        'content-type': 'application/json'
      },
      method: 'POST',
      mode: 'cors',
      redirect: 'follow',
      referrer: 'no-referrer',
    })
    .then( response => {
      return response.json()
    })
    .then( data => {
      if (data.code) {
        this.setState({ code: data.code })
      } 
    })
  }
  
	updateCode = ( newCode ) => {
		this.setState({
			code: newCode,
		})
	}

  render() {
    const {
      match
    } = this.props

    let options = {
      lineNumbers: true,
      readOnly: true,
      mode: 'php'
    }
    
    let styles = {
      display: 'block',
    }

    return (
      <div style={styles} className="reveal-overlay">
        <div style={styles} className="reveal full">
          {this.props.children}
          <div className="reveal-title"><h1>{match.file}</h1></div>
          <div className="reveal-content">
            <Highlighter value={this.state.code} file={match.file} options={options} />
          </div>
        </div>
      </div>
		);
  }
}

export default Modal