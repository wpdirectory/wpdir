import React, { Component } from 'react'

class ProgressBar extends Component {
  render() {
    const { progress } = this.props;
    let style = {
      width: progress+'%',
    }
    return (
	    <div className="progress">
        <div className="bar" style={style}></div>
      </div>
		);
  }
}

export default ProgressBar