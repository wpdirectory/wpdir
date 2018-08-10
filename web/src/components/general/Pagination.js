import React, { Component } from 'react'

class Pagination extends Component {
  render() {
    const {
      currentPage,
      totalPages,
    } = this.props
    
    let prev = ( currentPage <= 1 ) ? true : false
    let next = ( currentPage >= totalPages ) ? true : false
    let pageText = 'Page ' + currentPage + '/' + Math.ceil(totalPages)

    return (
	    <nav className="search-nav grid-x" aria-label="Pagination">
        <div className="cell small-6 medium-4"><button type="button" className="button primary hollow expanded" onClick={this.props.prevClick} disabled={prev}>Previous</button></div>
        <div className="cell small-6 medium-4 show-for-medium"><div className="search-nav-page">{pageText}</div></div>
        <div className="cell small-6 medium-4"><button type="button" className="button primary hollow expanded" onClick={this.props.nextClick} disabled={next}>Next</button></div>
      </nav>
	  )
  }
}

export default Pagination