import React, { Component } from 'react'
import Matches from './Matches.js'
//import Dashicon from '../Dashicon.js'

class Item extends Component {

  constructor(props) {
    super(props)
    this.state = {
      isActive: false,
    }
    this.toggleClass = this.toggleClass.bind(this)
  }

  formatName = (item) => {
    if (!item.name) {
      return item.slug
    } else {
      return item.name
    }
  }

  formatClass = () => {
    if (this.state.isActive === false) {
      return 'item'
    } else {
      return 'item active'
    }
  }

  toggleClass = () => {
    this.setState( prevState => ({
      isActive: !prevState.isActive
    }))
  }

  render() {
    const {
      id,
      item,
    } = this.props

    let matches
    if (this.state.isActive === true) {
      matches = <Matches id={id} slug={item.slug} />
    } else {
      matches = ''
    }

    return (
      <li className={this.formatClass()}>
        <button className="title" onClick={this.toggleClass}>
          <span className="name" dangerouslySetInnerHTML={{__html: this.formatName(item)}}></span>
          <span className="installs">{item.installs}</span>
          <span className="matches">{item.matches}</span>
        </button>
        <div className="content">
          {matches}
        </div>
      </li>
	  );
  }
}

export default Item