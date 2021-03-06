package main

import (
	"net"
	"net/http"

	"github.com/qaisjp/gosign"

	"github.com/gin-gonic/gin"
)

// Check does a thing whilsting checking something
// - success: abuse checks pass, and cosign returns a code + message
// - failure: cosign died somewhere, or the IP address is flagged for abuse
func (i *API) Check(c *gin.Context) {
	tokenName := c.Param("token_name")
	tokenKey := c.Param("token_key")
	if key, ok := i.Tokens[tokenName]; !ok || key != tokenKey {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	ip, ok := c.GetQuery("ip")
	if !ok || net.ParseIP(ip) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid IP provided",
		})
		return
	}

	cookie, ok := c.GetQuery("cookie")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "missing cookie",
		})
		return
	}

	response, err := i.GoSign.Check(cookie, false)

	if err == gosign.ErrLoggedOut {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "not logged in due to various reasons",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
