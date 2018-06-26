import React, { Component } from 'react'
import Item from './Item.js'
import Dashicon from '../Dashicon.js'

class Summary extends Component {

  constructor(props) {
    super(props)
    this.state = {
      id: this.props.id,
      items: this.props.items.sort( (a,b) => b.installs - a.installs ),
      desc: false,
      sorting: 'installs',
    }
    this.sortByName = this.sortByName.bind(this)
    this.sortByInstalls = this.sortByInstalls.bind(this)
    this.sortByMatches = this.sortByMatches.bind(this)
  }

  sortByName = () => {
    if (this.state.desc === true) {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'name',
        items: prevState.items.sort( (a,b) => {
          if (b.slug < a.slug) {
            return -1
          }
          if (b.slug > a.slug) {
            return 1
          }
          return 0
        })
      }))
    } else {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'name',
        items: prevState.items.sort( (a,b) => {
          if (a.slug < b.slug) {
            return -1
          }
          if (a.slug > b.slug) {
            return 1
          }
          return 0
        })
      }))
    }
  }

  sortByInstalls = () => {
    if (this.state.desc === true) {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'installs',
        items: prevState.items.sort( (a,b) => b.installs - a.installs )
      }))
    } else {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'installs',
        items: prevState.items.sort( (a,b) => a.installs - b.installs )
      }))
    }
  }

  sortByMatches = () => {
    if (this.state.desc === true) {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'matches',
        items: prevState.items.sort( (a,b) => b.matches - a.matches )
      }))
    } else {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'matches',
        items: prevState.items.sort( (a,b) => a.matches - b.matches )
      }))
    }
  }

  sortIcon = (name) => {
    if (name === this.state.sorting) {
      if (this.state.desc === false) {
        return <Dashicon icon="arrow-down-alt2" size={ 22 } />
      } else {
        return <Dashicon icon="arrow-up-alt2" size={ 22 } />
      }
    }
  }

  render() {
    let summaryItems
    if ( !!this.state.items && this.state.items.length && this.state.items.length > 0 ) {
      summaryItems = this.state.items.map( (item, key) => {
        return (
          <Item repo={this.props.repo} id={this.state.id} item={item} key={key} />
        );
      })
    } else {
      summaryItems = <p>Sorry, no matches found.</p>
    }

    return (
	    <ul className="accordion summary">
        <li className="header">
          <button className="name" onClick={this.sortByName}>Name{this.sortIcon('name')}</button>
          <button className="installs" onClick={this.sortByInstalls}>Installs{this.sortIcon('installs')}</button>
          <button className="matches" onClick={this.sortByMatches}>Matches{this.sortIcon('matches')}</button>
        </li>
        {summaryItems}
      </ul>
	  );
  }
}

export default Summary