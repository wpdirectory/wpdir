import React, { Component } from 'react'
import Matches from './Matches.js'

class Item extends Component {

  constructor(props) {
    super(props)
    this.state = {
      isActive: false,
    }
  }

  formatName = (item) => {
    if (!item.name) {
      return item.slug
    } else {
      return item.name
    }
  }

  formatInstalls = (item) => {
    if (!item.active_installs) {
      return 'n/a';
    }
    let installs = item.active_installs.toString()
    console.log(installs.length)
    if (installs.length > 6) {
      console.log('Over 6 Length')
      return installs.slice(0, installs.length - 6) + ',' + installs.slice(installs.length - 6, installs.length - 3) + ',' + installs.slice(installs.length - 3)
    }
    if (installs.length > 3) {
      return installs.slice(0, installs.length - 3) + ',' + installs.slice(installs.length - 3)
    }
    return installs
  }

  formatClass = () => {
    if (this.state.isActive === true) {
      return 'accordion-item is-active'
    } else {
      return 'accordion-item'
    }
  }

  toggleClass = () => {
    this.setState( prevState => ({
      isActive: !prevState.isActive
    }))
  }

  componentWillReceiveProps = (nextProps) => {
    if (nextProps.close === true) {
      this.setState({ isActive: false })
    }
  }

  render() {
    const {
      id,
      item,
    } = this.props

    let matches
    if (this.state.isActive === true) {
      matches = <Matches repo={this.props.repo} id={id} slug={item.slug} />
    } else {
      matches = ''
    }

    return (
      <li className={this.formatClass()}>
        <button className="accordion-title" onClick={this.toggleClass}>
          <span className="name" dangerouslySetInnerHTML={{__html: this.formatName(item)}}></span>
          <span className="installs" dangerouslySetInnerHTML={{__html: this.formatInstalls(item)}}></span>
          <span className="matches">{item.matches}</span>
        </button>
        <div className="accordion-content">
          {matches}
        </div>
      </li>
	  );
  }
}

export default Item