import React, { Component } from 'react'
import Dashicon from './Dashicon.js'

class Modal extends Component {
  render() {
    const { progress } = this.props;
    return (
      <div className="overlay">
        <div className="modal">
          <div className="close"><Dashicon icon="no" size={ 22 } /></div>
          <div className="content"></div>
        </div>
      </div>
		);
  }
}

export default Modal