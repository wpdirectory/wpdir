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
  }

  static defaultProps = {
    className: '',
    prefix: 'cm-',
  }

  componentWillMount = () => {
    this.calculateLanguage(this.props.file)
  }

  render() {
    const { value, prefix, lang } = this.props
  
    const elements = []
    let index = 0
    let lastStyle = null
    let tokenBuf = ''
    const pushElement = (token, style) => {
      elements.push(<span className={style ? prefix + style : ''} key={++index}>{token}</span>)
    }
    const mode = CodeMirror.findModeByName(lang)
    CodeMirror.runMode(value, mode ? mode.mime : 'meta', (token, style) => {
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
      <pre className={'cm-s-default'}>{elements}</pre>
	  )
  }
}

export default Highlighter