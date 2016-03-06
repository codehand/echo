package middleware

import (
	"net"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/webx-top/echo"
)

func Log() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger := c.Logger()

			remoteAddr := req.RemoteAddress()
			if ip := req.Header().Get(echo.XRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header().Get(echo.XForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method()
			path := req.URL().Path()
			if path == "" {
				path = "/"
			}
			size := res.Size()

			n := res.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			logger.Infof("%s %s %s %s %s %d", remoteAddr, method, path, code, stop.Sub(start), size)
			return nil
		})
	}
}