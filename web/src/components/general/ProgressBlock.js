import React, { Component } from 'react'

const messages = [
  'Filling Up Loading Bar',
  'Capitalizing P\'s',
  'Making Up Results',
  'Updating to PHP7',
  'Resolving Dependencies',
  'Listening To Jazz',
  'Initiating the Loop',
  'Yodaing conditions, I am',
  'do_action( \'search-directory\' )',
  'Hopefully not _doing_it_wrong()',
  'Pressing this',
  'Getting some REST'
]

class ProgressBlock extends Component {
  constructor(props) {
    super(props)
    this.state = {
      currentMessage: 0,
      timer: null,
      messageTimer: null,
      message: messages[Math.floor( Math.random() * messages.length )],
      counter: 0,
    }
  }

  componentDidMount = () => {
    let timer = setInterval( this.tick, 1000 )
    let messageTimer = setInterval( this.changeMessage, 5000 )
    let taken = Math.ceil( ( Date.now() - this.props.started ) / 1000 )
    this.setState({
      timer,
      messageTimer,
      counter: taken,
    })
  }

  componentWillUnmount = () => {
    clearInterval( this.state.timer )
    clearInterval( this.state.messageTimer )
  }

  tick = () => {
    this.setState({ counter: this.state.counter + 1 })
  }

  changeMessage = () => {
    this.setState({ message: messages[Math.floor( Math.random() * messages.length )] })
  }

  formatCounter = () => {
    var minutes = ( ( this.state.counter / 60 ) | 0 ) + ''
		var seconds = ( this.state.counter % 60 ) + ''
		return new Array( 3-minutes.length ).join('0') + minutes + ':' + new Array( 3-seconds.length ).join('0') + seconds
  }

  render() {
    const { progress } = this.props
    const { message } = this.state
    let style = {
      width: progress+'%',
    }
    return (
      <div className="progress panel cell small-12">
        <div className="progress-block">
          <div className="progress-block-background"></div>
          <div className="progress-block-bar" style={style}></div>
          <div className="progress-block-messages">{message}</div>
          <div className="progress-block-timer">{this.formatCounter()}</div>
        </div>
      </div>
		);
  }
}

export default ProgressBlock