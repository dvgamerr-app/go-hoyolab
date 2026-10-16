// Loads the Traceway browser SDK only when a DSN is configured; never breaks the page.
;(function () {
  try {
    var dsn = window.PUBLIC_TRACEWAY_DSN
    if (!dsn) return
    var s = document.createElement('script')
    s.src = 'https://cdn.jsdelivr.net/npm/@tracewayapp/frontend@1.2.1/dist/traceway.iife.global.js'
    s.integrity = 'sha384-QVzru2sR5BEVzzBnH9wcv8Pf3XhlbyCQ9Sjtknt6kJvjJJa/Q0XYInAOzTLogL3p'
    s.crossOrigin = 'anonymous'
    s.async = true
    s.onload = function () {
      try {
        window.Traceway.init(dsn, { recordAllSessions: true })
      } catch (e) {
        console.warn('Traceway disabled:', e && e.message ? e.message : e)
      }
    }
    document.head.appendChild(s)
  } catch (e) {
    console.warn('Traceway disabled:', e && e.message ? e.message : e)
  }
})()
