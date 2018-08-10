let Config = {
    Hostname: (window.config.Hostname === '%HOSTNAME%') ? 'https://wpdirectory.net' : window.config.Hostname,
    HTTP: {
        Timeout: (window.config.Timeout === '%TIMEOUT%') ? 5000 : window.config.Timeout,
    },
}

export default Config