import React, { Component } from 'react'

const messages = [
  'Filling Up Loading Bar',
  'Capitalizing P\'s',
  'Making Up Results',
  'Updating to PHP7',
  'Resolving Dependencies',
  'Listening To Jazz'
]

class ProgressBlock extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentMessage: 0,
      timer: null,
      counter: 0,
    }
  }

  componentDidMount = () => {
    let timer = setInterval( this.tick, 1000 )
    this.setState({ timer })
  }

  componentWillUnmount = () => {
    clearInterval( this.state.timer )
  }

  tick = () => {
    this.setState({ counter: this.state.counter + 1 })
  }

  formatCounter = () => {
    var minutes = ( ( this.state.counter / 60 ) | 0 ) + ''
		var seconds = ( this.state.counter % 60 ) + ''
		return new Array( 3-minutes.length ).join('0') + minutes + ':' + new Array( 3-seconds.length ).join('0') + seconds
  }

  render() {
    const { progress, status } = this.props;
    let style = {
      width: progress+'%',
    }
    return (
      <div className="progress panel cell small-12">
        <div className="progress-block">
          <div className="progress-block-background"></div>
          <div className="progress-block-bar" style={style}></div>
          <h2 className="progress-block-title">{status} - {this.formatCounter()}</h2>
        </div>
      </div>
		);
  }
}

export default ProgressBlock