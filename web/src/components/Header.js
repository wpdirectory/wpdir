import React, { Component } from 'react';
import { NavLink, Link } from 'react-router-dom';

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
          <Link to="/">WPDir</Link>
        </div>
        <a className={ this.state.toggle ? 'menu-toggle toggle active' : 'menu-toggle toggle' } aria-controls="primary-menu" aria-expanded="false" onClick={this.toggle}>
          <div className="toggle-box">
            <div className="toggle-inner"></div>
          </div>
        </a>

        <nav className={ this.state.toggle ? 'main-menu active' : 'main-menu' } aria-label="Main Navigation">
          <ul className="menu">
            <li><NavLink to="/searches">Searches</NavLink></li>
            <li><NavLink to="/stats">Stats</NavLink></li>
            <li><NavLink to="/reports">Reports</NavLink></li>
            <li><NavLink to="/about">About</NavLink></li>
          </ul>
          <div className="contact">
          <p>
            Want to talk? Send me a message:<br />
            <a href="mailto: mail@wpdirectory.net">mail@wpdirectory.net</a>
          </p>
          <ul className="social-links">
            <li><a href="https://www.twitter.com/wpdir/" target="_blank" rel="noopener noreferrer">Twitter</a></li>
            <li><a href="https://www.reddit.com/user/peterbooker/" target="_blank" rel="noopener noreferrer">Reddit</a></li>
            <li><a href="https://github.com/wpdirectory/wpdir" target="_blank" rel="noopener noreferrer">Github</a></li>
          </ul>
          </div>
        </nav>
      </header>
    )
  }
}

export default Header;