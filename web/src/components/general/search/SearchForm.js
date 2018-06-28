import React, { Component } from 'react'
import { withRouter } from 'react-router-dom'
import Dashicon from '../Dashicon.js'

class SearchForm extends Component {
    constructor(props) {
      super(props)
      this.state = {
          input: '',
          target: 'plugins',
          private: false
      }
    }

    updateInput = (event) => {
      this.setState({input: event.target.value})
    }

    updateTarget = (event) => {
      event.preventDefault();
      this.setState({target: event.target.value})
    }

    togglePrivate = () => {
      this.setState((prevState) => ({
        private: !prevState.private,
      }))
    }
  
    handleSubmit = (event) => {
      event.preventDefault();
      console.log(this.state)
      this.postData('https://wpdirectory.net/api/v1/search/new', {input: this.state.input, target: this.state.target, private: this.state.private})
        .then(data => {
            console.log(data)
            this.props.history.push('/search/' + data.id)
        })
        .catch(error => console.error(error))
    }

    postData(url, data) {
        return fetch(url, {
          body: JSON.stringify(data),
          cache: 'no-cache',
          credentials: 'same-origin',
          headers: {
            'user-agent': 'WPDirectory/0.1.0',
            'content-type': 'application/json'
          },
          method: 'POST',
          mode: 'cors',
          redirect: 'follow',
          referrer: 'no-referrer',
        })
        .then(response => response.json())
    }
  
    render() {
      return (
        <form className="search-form" onSubmit={this.handleSubmit}>
          <h3>New Search</h3>
          <div className="input-choice">
            <label>Regular Expression:</label>
            <input className="input" type="text" placeholder="" value={this.state.input} onChange={this.updateInput} />
          </div>

          <div className="directory-choice">
            <label>What to Search:</label>
            <div className="button-group expanded stacked-for-small">
              <button className={this.state.target === 'plugins' ? 'button' : 'button secondary'} value="plugins" onClick={this.updateTarget}>
                <Dashicon icon="admin-plugins" size={ 16 } />
                Plugins
              </button>
              <button className={this.state.target === 'themes' ? 'button' : 'button secondary'} value="themes" onClick={this.updateTarget} disabled={true}>
                <Dashicon icon="admin-appearance" size={ 16 } />
                Themes
              </button>
            </div>
          </div>

          <div className="private-choice">
            <label>Make search private?</label>
            <div className="switch large">
              <input className="switch-input" id="yes-no" type="checkbox" name="privateSwitch" defaultChecked={false} />
              <label className="switch-paddle" htmlFor="yes-no" onClick={this.togglePrivate}>
                <span className="show-for-sr">Make search private?</span>
                <span className="switch-active" aria-hidden="true">Yes</span>
                <span className="switch-inactive" aria-hidden="true">No</span>
              </label>
            </div>
          </div>

          <input className="button expanded" type="submit" value="Search" onClick={this.handleSubmit} />
        </form>
      );
    }
}

export default withRouter(SearchForm)