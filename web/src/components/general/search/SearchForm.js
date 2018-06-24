import React, { Component } from 'react'
import { withRouter } from 'react-router-dom';
import Dashicon from '../Dashicon.js'

class SearchForm extends Component {
    constructor(props) {
      super(props);
      this.state = {
          input: '',
          target: ''
      };
  
      this.handleChange = this.handleChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
    }
  
    handleChange(event) {
      this.setState({input: event.target.value});
    }
  
    handleSubmit(event) {
      event.preventDefault();
      this.postData('https://wpdirectory.net/api/v1/search/new', {input: this.state.input, target: 'plugins'})
        .then(data => {
            console.log(data)
            this.props.history.push('/search/' + data.id)
        }) // JSON from `response.json()` call
        .catch(error => console.error(error))
    }

    postData(url, data) {
        // Default options are marked with *
        return fetch(url, {
          body: JSON.stringify(data), // must match 'Content-Type' header
          cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
          credentials: 'same-origin', // include, same-origin, *omit
          headers: {
            'user-agent': 'WPDirectory/0.1.0',
            'content-type': 'application/json'
          },
          method: 'POST', // *GET, POST, PUT, DELETE, etc.
          mode: 'cors', // no-cors, cors, *same-origin
          redirect: 'follow', // manual, *follow, error
          referrer: 'no-referrer', // *client, no-referrer
        })
        .then(response => response.json()) // parses response to JSON
    }
  
    render() {
      return (
        <form className="search-form" onSubmit={this.handleSubmit}>
          <h3>New Search</h3>
          <div className="input-choice">
            <label>Regular Expression:</label>
            <input className="input" type="text" placeholder="" value={this.state.input} onChange={this.handleChange} />
          </div>
          <div className="directory-choice">
            <label>What to Search:</label>
            <div className="button-group">
              <input type="radio" id="search-plugins" name="target" value="Plugins" defaultChecked={true} />
              <label className="button" htmlFor="search-plugins">
                <Dashicon icon="admin-plugins" size={ 22 } />
                Plugins
              </label>
              <input type="radio" id="search-themes" name="target" value="Themes" />
              <label className="button" htmlFor="search-themes">
                <Dashicon icon="admin-appearance" size={ 22 } />
                Themes
              </label>
            </div>
          </div>
          <input type="submit" value="Search" />
        </form>
      );
    }
}

export default withRouter(SearchForm)