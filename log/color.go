package log

import log "github.com/corgi-kx/logcustom"

type ColorLog struct {
}

func (c *ColorLog) Write(p []byte) (n int, err error) {
	log.Trace(string(p))
	return
}

func (c *ColorLog) Sync() (err error) {
	return
}
