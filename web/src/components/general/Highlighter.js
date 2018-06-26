import React, { Component } from 'react'
import CodeMirror from 'codemirror'
import 'codemirror/addon/runmode/runmode'
import 'codemirror/mode/meta'
import 'codemirror/mode/php/php'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/mode/css/css'
import 'codemirror/mode/sql/sql'
import 'codemirror/mode/markdown/markdown'
import 'codemirror/mode/htmlmixed/htmlmixed'

import 'codemirror/lib/codemirror.css'

class Highlighter extends Component {

  constructor(props) {
    super(props)
    this.state = {
      language: 'meta'
    }
  }

  static defaultProps = {
    className: '',
    prefix: 'cm-',
    theme: 'default'
  }

  calculateLanguage = (filename) => {
    let guess = filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2)
    switch (guess) {
      case 'php':
        this.setState({language: 'php'});
        break
      case 'js':
        this.setState({language: 'javascript'});
        break
      case 'css':
        this.setState({language: 'css'});
        break
      case 'sass':
      case 'scss':
        this.setState({language: 'sass'});
        break
      case 'sql':
        this.setState({language: 'sql'});
        break
      case 'md':
        this.setState({language: 'markdown'});
        break
      case 'htm':
      case 'html':
        this.setState({language: 'htmlmixed'});
        break
      default:
        this.setState({language: 'meta'});
    }
  }

  componentWillMount = () => {
    this.calculateLanguage(this.props.file)
  }

  render() {
    const { value, file, prefix } = this.props
  
    console.log(this.state.language)
    const elements = []
    let index = 0
    let lastStyle = null
    let tokenBuf = ''
    const pushElement = (token, style) => {
      elements.push(<span className={style ? prefix + style : ''} key={++index}>{token}</span>)
    }
    const mode = CodeMirror.findModeByName(this.state.language)
    CodeMirror.runMode(value, mode ? mode.mime : 'meta', (token, style) => {
      //console.log(token)
      //console.log(style)
      if (lastStyle === style) {
        tokenBuf += token
        lastStyle = style
      } else {
        if (tokenBuf) {
          pushElement(tokenBuf, lastStyle)
          
        }
        tokenBuf = token
        lastStyle = style
      }
    })
    pushElement(tokenBuf, lastStyle)

    return (
		  <div className="codemirror">
        <pre className={'cm-s-default'}>{elements}</pre>
      </div>
	  );
  }
}

export default Highlighter