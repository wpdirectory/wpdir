import React, { Component } from 'react'
import Dashicon from '../Dashicon.js'
import Modal from './Modal.js'

const ESCAPE = 27

class Match extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modalActive: false,
    }
  }

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

  toggleModal = () => {
    this.setState((prevState) => ({
      modalActive: !prevState.modalActive,
    }))
  }

  escFunction = (event) => {
    if ( event.keyCode === ESCAPE ) {
      this.setState({ modalActive: false })
    }
  }

  componentDidMount = () => {
    document.addEventListener( 'keydown', this.escFunction, false )
  }

  componentWillUnmount = () => {
    document.removeEventListener( 'keydown', this.escFunction, false )
  }
  
  render() {
    const {
      match,
      repo,
    } = this.props

    let modal
    if (this.state.modalActive) {
      modal = (
        <Modal isOpen={this.state.modalActive} repo={repo} match={match}>
          <button className="close-button" onClick={this.toggleModal}>
            <span aria-hidden="true">&times;</span>
          </button>
        </Modal>
      )
    }

    return (
      <li className="match">
        <span className="num">{match.line_num}</span>
        <span className="text"><code>{match.line_text}</code></span>
        <button className="view" onClick={this.toggleModal}><Dashicon icon="editor-code" size={ 22 } /></button>
        {modal}
      </li>
	  )
  }
}

export default Match