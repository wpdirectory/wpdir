import React, { Component } from 'react';

class NotFound extends Component {

  componentDidMount(){
    document.title = 'Not Found - WPdirectory'
  }

  render(){
      return (
        <div className="page page-404">
          <div className="grid-x gutter-x gutter-y small-1">
            <div className="panel">
              <h1>404 - Not Found</h1>
              <p>Oops, we could not find what you were looking for.</p>
          </div>
          </div>
        </div>
      );
  }
}

export default NotFound;