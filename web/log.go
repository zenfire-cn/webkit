package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func setGinLog(logger *zap.Logger, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("S", c.Writer.Status()),
				zap.String("M", c.Request.Method),
				zap.Duration("T", latency),
				zap.String("Q", query),
				zap.String("IP", c.ClientIP()),
				// zap.String("UA", c.Request.UserAgent()),
			)
		}
	}
}

func setGinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				path := c.Request.URL.Path
				errInfo := path + "\t " + err.(error).Error()
				if brokenPipe {
					logger.Error(errInfo)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error))
					c.Abort()
					return
				}
				if stack {
					logger.Error(errInfo+"\n"+getStack(), zap.String("REQUEST", string(httpRequest)))
				} else {
					logger.Error(errInfo)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func getStack() string {
	stacks := strings.Split(string(debug.Stack()), "\n")[9:]
	return strings.Join(stacks, "\n")
}
