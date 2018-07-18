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
          <span className="installs">{item.active_installs}</span>
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