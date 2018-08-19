import React, { Component } from 'react'
import Dashicon from '../Dashicon.js'
import Matches from './Matches.js'

class Item extends Component {
  constructor(props) {
    super(props)
    this.state = {
      isActive: false,
    }
  }

  formatName = ( item ) => {
    if ( !item.name ) {
      return item.slug
    } else {
      return item.name
    }
  }

  formatInstalls = (item) => {
    if (!item.active_installs) {
      return '0';
    }
    let installs = item.active_installs.toString()
    if (installs.length > 6) {
      return installs.slice( 0, installs.length - 6 ) + ',' + installs.slice( installs.length - 6, installs.length - 3 ) + ',' + installs.slice (installs.length - 3 )
    }
    if (installs.length > 3) {
      return installs.slice( 0, installs.length - 3 ) + ',' + installs.slice( installs.length - 3 )
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

  stopPropagation = (e) => {
    e.stopPropagation()
  }

  render() {
    const {
      id,
      item,
      repo
    } = this.props

    let matches
    if (this.state.isActive === true) {
      matches = <Matches repo={this.props.repo} id={id} slug={item.slug} />
    } else {
      matches = ''
    }

    let title = ( item.name ? item.name : item.slug )

    return (
      <li className={this.formatClass()}>
        <div className="accordion-title" onClick={this.toggleClass}>
          <span className="name">
            <a className="wplink" href={'https://wordpress.org/' + repo + '/' + item.slug + '/'} title={'WordPress.org - ' + title} onClick={this.stopPropagation} target="_blank" rel="noopener noreferrer">
              <Dashicon icon="wordpress" size={ 22 } />
            </a>
            <span dangerouslySetInnerHTML={{__html: this.formatName(item)}}></span>
          </span>
          <span className="installs">{this.formatInstalls(item)}</span>
          <span className="matches">{item.matches}</span>
        </div>
        <div className="accordion-content">
          {matches}
        </div>
      </li>
	  );
  }
}

export default Item