import React from 'react';
import { BrowserRouter as Link } from 'react-router-dom';

const Header = () => {
  return (
    <header className="header grid-x gutter-x">
      <div className="title">
        <Link to="/">mtStats</Link>
      </div>
      <nav className="main-menu" aria-label="Main Navigation">
      <ul className="menu">
        <li><Link to="/news">News</Link></li>
        <li className="has-dropdown">
          <Link to="/data">Data <i className="fas fa-chevron-down"></i></Link>
          <ul className="sub-menu">
            <li><Link to="/news">Fights</Link></li>
            <li><Link to="/news">Fighters</Link></li>
            <li><Link to="/news">Events</Link></li>
            <li><Link to="/news">Organisations</Link></li>
          </ul>
        </li>
        <li><Link to="/data">Fights</Link></li>
        <li><Link to="/data">About</Link></li>
        <li><Link to="/join">Join <i className="fas fa-user"></i></Link></li>
      </ul>
      </nav>
    </header>
  );
};

export default Header;