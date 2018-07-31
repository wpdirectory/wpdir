import React, { Component } from 'react'
import Item from './Item.js'
import Dashicon from '../Dashicon.js'

class Summary extends Component {

  constructor(props) {
    super(props)
    this.state = {
      id: this.props.id,
      matches: this.props.matches,
      items: [],
      sorting: 'installs',
      desc: true,
      currentPage: 1,
      perPage: 250,
      isLoading: true,
    }
  }

  componentWillMount = () => {
    fetch(window.wpdirHost + '/api/v1/search/summary/' + this.props.id)
    .then( response => {
      return response.json()
    })
    .then( data => {
      this.setState({ items: data.results })
      this.setState({ isLoading: false })
      this.sortByInstalls()
    })
  }

  sortByName = () => {
    this.setState({ forceClose: true })
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
    this.setState({ forceClose: true })
    if (this.state.desc === true) {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'installs',
        items: prevState.items.sort( (a,b) => b.active_installs - a.active_installs )
      }))
    } else {
      this.setState((prevState) => ({
        desc: !prevState.desc,
        sorting: 'installs',
        items: prevState.items.sort( (a,b) => a.active_installs - b.active_installs )
      }))
    }
  }

  sortByMatches = () => {
    this.setState({ forceClose: true })
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

  prevPage = () => {
    this.setState((prevState) => ({
      currentPage: prevState.currentPage - 1,
    }))
  }

  nextPage = () => {
    this.setState((prevState) => ({
      currentPage: prevState.currentPage + 1,
    }))
  }

  render() {
    const { 
      items,
      currentPage,
      perPage
    } = this.state

    console.log(currentPage)

    const indexOfLastItem = currentPage * perPage
    const indexOfFirstItem = indexOfLastItem - perPage
    const currentItems = items.slice(indexOfFirstItem, indexOfLastItem)
    const numPages = items.length / perPage

    let summaryItems
    if ( !!this.state.items && this.state.items.length && this.state.items.length > 0 ) {
      summaryItems = currentItems.map( (item, key) => {
        return (
          <Item repo={this.props.repo} id={this.state.id} item={item} close={this.state.forceClose} key={key} />
        )
      })
    } else {
      summaryItems = <p>Sorry, no matches found.</p>
    }

    return (
      <div>
        <ul className="accordion summary">
          <li className="header">
            <button className="name" onClick={this.sortByName}>Name{this.sortIcon('name')}</button>
            <button className="installs" onClick={this.sortByInstalls}>Installs{this.sortIcon('installs')}</button>
            <button className="matches" onClick={this.sortByMatches}>Matches{this.sortIcon('matches')}</button>
          </li>
          {summaryItems}
        </ul>
        <nav aria-label="Pagination">
          <ul className="pagination text-center">
            <li className="pagination-previous"><button aria-label="Previous page" onClick={this.prevPage}>Previous <span className="show-for-sr">page</span></button></li>
            <li className="pagination-next"><button aria-label="Next page" onClick={this.nextPage}>Next <span className="show-for-sr">page</span></button></li>
          </ul>
        </nav>
      </div>
	  );
  }
}

export default Summary