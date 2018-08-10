import React, { Component } from 'react'
import Item from './Item.js'
import Dashicon from '../Dashicon.js'
import Loadicon from '../Loadicon.js'
import Pagination from '../Pagination.js'
import API from '../../../utils/API.js'

class Summary extends Component {

  constructor(props) {
    super(props)
    this.state = {
      id: this.props.id,
      matches: this.props.matches,
      items: [],
      sorting: 'installs',
      desc: false,
      currentPage: 1,
      perPage: 100,
      isLoading: true,
      error: ''
    }
    this.topRef = React.createRef()
  }

  componentWillMount = () => {
    this.setState({ isLoading: true })

    API.get( '/search/summary/' + this.props.id )
      .then( result => {
        this.setState({
          items: result.data.results,
          isLoading: false
        })
        this.sortByInstalls()
      })
      .catch( error => this.setState({
        error,
        isLoading: false
      }))

    this.sortByInstalls()
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
    this.topRef.current.scrollIntoView({
      block: "start",
      behavior: "smooth"
    })
  }

  nextPage = () => {
    this.setState((prevState) => ({
      currentPage: prevState.currentPage + 1,
    }))
    this.topRef.current.scrollIntoView({
      block: "start",
      behavior: "smooth"
    })
  }

  render() {
    const { 
      items,
      currentPage,
      perPage,
      isLoading,
      error
    } = this.state

    const indexOfLastItem = currentPage * perPage
    const indexOfFirstItem = indexOfLastItem - perPage
    const currentItems = items.slice(indexOfFirstItem, indexOfLastItem)
    const numPages = items.length / perPage
    const needsPagination = ( numPages <= 1 ) ? false : true

    let summaryItems

    if ( isLoading ) {
      return (
        <div className="search-summary panel cell small-12">
          <Loadicon />
        </div>
      )
    }  else {
      if ( error ) {
        summaryItems = <p className="error">Sorry, there was a problem fetching data.</p>
      } else {
        summaryItems = currentItems.map( (item, key) => {
          return (
            <Item repo={this.props.repo} id={this.state.id} item={item} close={this.state.forceClose} key={key} />
          )
        })
      }
    }

    return (
      <div className="search-summary panel cell small-12">
        <h2>Summary <small>({items.length}{ ' Extensions'})</small></h2>
        <div ref={this.topRef}>
          { needsPagination && <Pagination currentPage={currentPage} totalPages={numPages} prevClick={this.prevPage} nextClick={this.nextPage} /> }
          <ul className="accordion summary">
            <li className="header">
              <button className="name" onClick={this.sortByName}>Name{this.sortIcon('name')}</button>
              <button className="installs" onClick={this.sortByInstalls}>Installs{this.sortIcon('installs')}</button>
              <button className="matches" onClick={this.sortByMatches}>Matches{this.sortIcon('matches')}</button>
            </li>
            { summaryItems }
          </ul>
          { needsPagination && <Pagination currentPage={currentPage} totalPages={numPages} prevClick={this.prevPage} nextClick={this.nextPage} /> }
        </div>
      </div>
	  )
  }
}

export default Summary