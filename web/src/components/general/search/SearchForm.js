import React, { Component } from 'react'
import { withRouter } from 'react-router-dom'
import Dashicon from '../Dashicon.js'
import API from '../../../utils/API.js'

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
      event.preventDefault()
      this.setState({
        isLoading: true,
        error: ''
      })
      API.post( '/search/new', {
        input: this.state.input,
        target: this.state.target,
        private: this.state.private
      })
      .then( response => {
        this.setState({
          isLoading: false
        })
        this.props.history.push( '/search/' + response.data.id )
      })
      .catch( error => {
        if ( error.response ) {
          this.setState({
            error: 'Error bad status: ' + error.response.status,
            isLoading: false
          })
        } else if ( error.request ) {
          this.setState({
            error: 'Error no response received',
            isLoading: false
          })
        } else {
          this.setState({
            error: 'Error making request: ' + error.message,
            isLoading: false
          })
        }
      })
    }

    render() {
      const { 
        input,
        target,
        isLoading,
        error
      } = this.state

      return (
        <form className="search-form" onSubmit={this.handleSubmit}>
          <h2>New Search</h2>
          <div className="input-choice">
            <label>Regular Expression:</label>
            <input className="input" type="text" placeholder="" value={input} onChange={this.updateInput} />
          </div>

          <div className="directory-choice">
            <label>What to Search:</label>
            <div className="button-group expanded stacked-for-small">
              <button className={target === 'plugins' ? 'button' : 'button secondary'} value="plugins" onClick={this.updateTarget}>
                <Dashicon icon="admin-plugins" size={ 16 } />
                Plugins
              </button>
              <button className={target === 'themes' ? 'button' : 'button secondary'} value="themes" onClick={this.updateTarget}>
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

          { ( !isLoading && error ) &&
          <div className="callout alert">
            <p>{error}</p>
          </div>
          }
        </form>
      );
    }
}

export default withRouter(SearchForm)