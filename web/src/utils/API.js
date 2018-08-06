import axios from 'axios'
import Hostname from './Hostname.js'

let API = axios.create({
  baseURL: Hostname + '/api/v1',
  timeout: 5000,
})

export default API