import React, { Component } from 'react';
import { NavLink, Link } from 'react-router-dom';
import Logo from './general/logos/Logo.js'

class Header extends Component {
  constructor(props) {
      super(props);
      this.state = {
          toggle: false
      }
  }

  toggle = () => {
    this.setState({ toggle: ! this.state.toggle })
  }
  
  render () {
    return (
      <header className="header cell shrink medium-cell-block-container">
        <div className="grid-x grid-margin-x">

          <div className="title cell small-6 medium-4 large-3">
            <Link to="/"><Logo width={255} height={50} /></Link>
          </div>

          <div className="toggle-area cell small-6 medium-8 large-9">
            <a className={ this.state.toggle ? 'toggle active' : 'toggle' } aria-controls="primary-menu" aria-expanded="false" onClick={this.toggle}>
              <div className="toggle-box">
                <div className="toggle-inner"></div>
              </div>
            </a>
          </div>

          <nav className={ this.state.toggle ? 'main-menu active cell small-12 medium-8 large-9' : 'main-menu cell small-12 medium-8 large-9' } aria-label="Main Navigation">
            <ul className="menu vertical medium-horizontal align-right">
              <li><NavLink onClick={this.toggle} exact to="/">Home</NavLink></li>
              <li><NavLink onClick={this.toggle} to="/searches">Searches</NavLink></li>
              <li><NavLink onClick={this.toggle} to="/repos">Repos</NavLink></li>
              <li><NavLink onClick={this.toggle} to="/about">About</NavLink></li>
            </ul>
            <div className="contact show-for-small-only">
              <p>Feedback or questions? We would love to hear from you!</p>
              <ul className="links">
                <li><a href="mailto: mail@peterbooker.com">Email</a></li>
                <li><a href="https://www.twitter.com/peter_booker/" target="_blank" rel="noopener noreferrer">Twitter</a></li>
                <li><a href="https://github.com/wpdirectory/wpdir" target="_blank" rel="noopener noreferrer">Github</a></li>
              </ul>
            </div>
          </nav>
        </div>
      
      </header>
    )
  }
}

export default Header;