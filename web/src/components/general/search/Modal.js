import React, { Component } from 'react'
import CodeMirror from '../CodeMirror.js'
import API from '../../../utils/API.js'

class Modal extends Component {
  constructor(props) {
    super(props)
    this.state = {
      code: '',
      lang: '',
      isLoading: true,
      error: ''
    }
  }

  getInitialState = () => {
		return {
			code: "// Code",
		}
  }

  componentWillMount = () => {
    this.setState({ isLoading: true })

    API.post( '/file', {
      repo: this.props.repo,
      slug: this.props.match.slug,
      file: this.props.match.file
      })
      .then( result => this.setState({
        code: result.data.code,
        isLoading: false
      }))
      .catch(error => this.setState({
        error,
        isLoading: false
      }))
  }

	updateCode = ( newCode ) => {
		this.setState({
			code: newCode,
		})
  }
  
  calculateLanguage = (filename) => {
    let guess = filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2)
    switch ( guess ) {
      case 'php':
        return 'php'
      case 'js':
        return 'javascript'
      case 'css':
        return 'css'
      case 'sass':
      case 'scss':
        return 'sass'
      case 'sql':
        return 'sql'
      case 'md':
        return 'markdown'
      case 'htm':
      case 'html':
        return 'htmlmixed'
      default:
        return 'meta'
    }
  }

  render() {
    const {
      match
    } = this.props

    const {
      isLoading,
      error
    } = this.state

    let options = {
      styleActiveLine: true,
      lineNumbers: true,
      readOnly: true,
      mode: this.calculateLanguage(match.file)
    }
    
    let styles = {
      display: 'block',
    }

    return (
      <div style={styles} className="reveal-overlay">
        <div style={styles} className="reveal large">
          <h1>{match.file}</h1>
          <div className="reveal-content">
            <CodeMirror value={this.state.code} file={match.file} line={match.line_num} options={options} focus={true} />
          </div>
          {this.props.children}
        </div>
      </div>
		);
  }
}

export default Modal