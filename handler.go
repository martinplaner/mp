package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type queryParams struct {
	Query string `form:"q,default=mp"`
}

func handlerWithGenerator(g *Generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params queryParams
		c.BindQuery(&params)

		if len(params.Query) < 1 || len(params.Query) > 20 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		word, err := g.Generate(strings.ToUpper(params.Query))
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEHTML, gin.MIMEJSON) {
		case binding.MIMEPlain:
			c.String(http.StatusOK, word)
		case binding.MIMEJSON:
			c.JSON(http.StatusOK, gin.H{"query": params.Query, "word": word})
		case binding.MIMEHTML:
			c.HTML(http.StatusOK, "index.tmpl", gin.H{"Word": word})
		default:
			c.AbortWithError(http.StatusNotAcceptable, errors.New("the target resource does not have a current representation that would be acceptable to the user agent"))
		}
	}
}
