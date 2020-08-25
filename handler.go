package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	// DefaultQuery is used if no other query was specified
	DefaultQuery = "mp"
)

// QueryResult contains the query as well as the result value of a single request
type QueryResult struct {
	Query  string `json:"query"`
	Result string `json:"result"`
}

// Error holds useful information in case of an error
type Error struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

func handlerWithGenerator(g *Generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := strings.TrimSpace(strings.TrimPrefix(c.Request.URL.Path, "/"))

		if strings.Contains(query, "/") || len(query) > 20 {
			c.AbortWithStatusJSON(http.StatusBadRequest, &Error{
				Message: "invalid query",
			})
			return
		} else if query == "" {
			query = DefaultQuery
		}

		result, err := g.Generate(strings.ToUpper(query))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &Error{
				Message: fmt.Sprintf("could not generate word for '%v' query", query),
				Cause:   err.Error(),
			})
			return
		}

		queryResult := &QueryResult{Query: query, Result: result}

		switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEHTML, gin.MIMEJSON) {
		case binding.MIMEPlain:
			c.String(http.StatusOK, queryResult.Result)
		case binding.MIMEJSON:
			c.JSON(http.StatusOK, queryResult)
		case binding.MIMEHTML:
			c.HTML(http.StatusOK, "index.tmpl", queryResult)
		default:
			c.AbortWithError(http.StatusNotAcceptable, errors.New("the target resource does not have a current representation that would be acceptable to the user agent"))
		}
	}
}
