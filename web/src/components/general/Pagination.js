import React, { Component } from 'react'

class Pagination extends Component {
  render() {
    const {
      currentPage,
      totalPages,
    } = this.props
    let prev = ( currentPage <= 1 ) ? true : false
    let next = ( currentPage >= totalPages ) ? true : false
    return (
	    <nav className="search-nav expanded button-group" aria-label="Pagination">
        <button type="button" className="button primary" onClick={this.props.prevClick} disabled={prev}>Previous</button>
        <button type="button" className="button primary" onClick={this.props.nextClick} disabled={next}>Next</button>
      </nav>
	  )
  }
}

export default Pagination