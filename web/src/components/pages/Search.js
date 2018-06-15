import React, { Component } from 'react'
import Dashicon from '../general/Dashicon.js'

class Search extends Component {
  render() {

    let data = [
      {
        "slug": "health-check",
        "file": "health-check/health-check.php",
        "lines": [
          {
            "number": "45",
            "text": "// Set the current cURL version."
          },
          {
            "number": "46",
            "text": "define( 'HEALTH_CHECK_CURL_VERSION', '7.58' );"
          },
          {
            "number": "47",
            "text": ""
          },
          {
            "number": "48",
            "text": "// Set the minimum cURL version that we've tested that core works with."
          },
          {
            "number": "49",
            "text": "define( 'HEALTH_CHECK_CURL_MIN_VERSION', '7.38' );"
          },
        ],
      },
      {
        "slug": "health-check",
        "file": "health-check/includes/class-health-check-auto-updates.php",
        "lines": [
          {
            "number": "39",
            "text": "     * @return array"
          },
          {
            "number": "40",
            "text": "     */"
          },
          {
            "number": "41",
            "text": "    public function run_tests() {"
          },
          {
            "number": "42",
            "text": "        $tests = array();"
          },
          {
            "number": "43",
            "text": ""
          },
        ],
      },
    ]

    let searchResults = data.map( function(result, idx) {
      return (
        <div key={idx} className="result">
          <div className="file">
            <span className="name">{result.file}</span>
            <a className="link" href={"https://plugins.trac.wordpress.org/browser/" + result.file} target="_blank" rel="noopener noreferrer">
              <Dashicon icon="external" size={ 22 } />
            </a>
          </div>
          <ul className="lines">
            {result.lines.map(function(line, idx) {
              return (
              <li key={idx}>
                <span className="num">{line.number}</span>
                <span className="excerpt"><code>{line.text}</code></span>
              </li>
              )
            })}
          </ul>
        </div>
      )
    })

      return (
        <div className="page page-search">
          <div className="title panel">
            <h1>Search</h1>
          </div>
          <div className="search-info panel">
            <h2>Overview</h2>
            <div className="info">
              <dl>
                <dt>Total Matches</dt>
                <dd>317</dd>
              </dl>
            </div>
          </div>
          <div className="search-results panel">
            <h2>Results</h2>
            <div className="results">
              {searchResults}
            </div>
          </div>
        </div>
      )
  }
}

export default Search