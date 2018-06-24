import React, { Component } from 'react'
import Dashicon from '../Dashicon.js'

class Match extends Component {
  //constructor(props) {
    //super(props)
  //}

  formatFilename = () => {
    let len = this.props.match.slug.length + 1
    return this.props.match.file.slice(len)
  }

  formatLine = () => {
    const {
      match,
    } = this.props
    return match.line_num + '' + match.line_text
  }
  
  render() {
    const {
      match,
    } = this.props

    return (
      <li className="match">
        <span className="num">{match.line_num}</span>
        <span className="text"><code>{match.line_text}</code></span>
        <button className="view"><Dashicon icon="editor-code" size={ 22 } /></button>
      </li>
	  );
  }
}

export default Match