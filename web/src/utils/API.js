import axios from 'axios'
import Config from './Config.js'

let API = axios.create({
  baseURL: Config.Hostname + '/api/v1',
  timeout: Config.HTTP.Timeout,
})

export default API