import React from 'react'
import Config from '../utils/Config.js'

const Footer = () => {
  return (
    <footer className="footer cell shrink">
      <div className="info">&copy; 2018 WP Directory - Made with Love, Go and React by <a href="https://www.peterbooker.com" target="_blank" rel="noopener noreferrer">Peter Booker</a> - <a href="https://github.com/wpdirectory/wpdir" target="_blank" rel="noopener noreferrer" title={'Version: v' + Config.Version + ' Commit: ' + Config.Commit + ' Date: ' + Config.Date}>{ 'v' + Config.Version }</a></div>
    </footer>
  )
}

export default Footer