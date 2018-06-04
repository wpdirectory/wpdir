import React, { Component } from 'react';

class Home extends Component {
  render(){
      return (
        <div className="page-home">
          <h1>Home</h1>
          <div className="grid-x gutter-x gutter-y">
            <div>One</div>
            <div>Two</div>
            <div>Three</div>
            <div>Four</div>
          </div>
        </div>
      );
  }
}

export default Home;