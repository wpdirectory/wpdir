let Config = {
    Version: (window.config.Version === '%VERSION%') ? '1.0.0' : window.config.Version,
    Commit: (window.config.Commit === '%COMMIT%') ? 'n/a' : window.config.Commit,
    Date: (window.config.Date === '%DATE%') ? '2000-01-01 00:00:00' : window.config.Date,
    Hostname: (window.config.Hostname === '%HOSTNAME%') ? 'https://wpdirectory.net' : window.config.Hostname,
    HTTP: {
        Timeout: (window.config.Timeout === '%TIMEOUT%') ? 5000 : window.config.Timeout,
    },
}

export default Config