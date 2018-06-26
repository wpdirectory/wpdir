import React, { Component } from 'react';

class NotFound extends Component {

  componentDidMount(){
    document.title = 'Not Found - WPdirectory'
  }

  render(){
      return (
        <div className="page page-404 grid-container">
          <div className="grid-x grid-margin-x grid-margin-y">
            <div className="panel cell small-12">
              <h1>404 - Not Found</h1>
              <p>Oops, we could not find what you were looking for.</p>
          </div>
          </div>
        </div>
      );
  }
}

export default NotFound;