import React, { Component } from 'react';
import { NavLink, Link } from 'react-router-dom';
//import LogoIcon from './general/logos/LogoIcon.js'
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
    console.log( this.state.toggle )
  }
  
  render () {
    return (
      <header className="header grid-x gutter-x">
        <div className="title">
          <Link to="/"><Logo width={255} height={50} /></Link>
        </div>
        <a className={ this.state.toggle ? 'menu-toggle toggle active' : 'menu-toggle toggle' } aria-controls="primary-menu" aria-expanded="false" onClick={this.toggle}>
          <div className="toggle-box">
            <div className="toggle-inner"></div>
          </div>
        </a>

        <nav className={ this.state.toggle ? 'main-menu active' : 'main-menu' } aria-label="Main Navigation">
          <ul className="menu">
            <li><NavLink to="/searches">Searches</NavLink></li>
            <li><NavLink to="/repos">Repos</NavLink></li>
            <li><NavLink to="/about">About</NavLink></li>
          </ul>
          <div className="contact">
          <p>Feedback or questions?</p>
          <ul className="social-links">
            <li><a href="mailto: mail@peterbooker.com">Email</a></li>
            <li><a href="https://www.twitter.com/peter_booker/" target="_blank" rel="noopener noreferrer">Twitter</a></li>
            <li><a href="https://github.com/wpdirectory/wpdir" target="_blank" rel="noopener noreferrer">Github</a></li>
          </ul>
          </div>
        </nav>
      </header>
    )
  }
}

export default Header;