import React from 'react'
import Config from '../utils/Config.js'
import DreamHostLink from './general/DreamHostLink.js'

const Footer = () => {
  return (
    <footer className="footer cell shrink">
      <div className="info">
        <span>Made with Love, Go and React by <a href="https://www.peterbooker.com" target="_blank" rel="noopener noreferrer">Peter Booker</a></span>&nbsp;-&nbsp;
        <span>Powered by <DreamHostLink height="16" width="120" /></span>&nbsp;-&nbsp;
        <span><a href="https://github.com/wpdirectory/wpdir" target="_blank" rel="noopener noreferrer" title={'Version: v' + Config.Version + ' Commit: ' + Config.Commit + ' Date: ' + Config.Date}>wpdir { Config.Version }</a></span>
      </div>
    </footer>
  )
}

export default Footer