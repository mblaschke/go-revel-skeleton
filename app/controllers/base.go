package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"net/http"
)

type Base struct {
	*revel.Controller
}

func (c Base) accessCheck() (result revel.Result) {
	return
}

func (c Base) renderJSONError(err string) (revel.Result) {
	c.Response.Status = http.StatusInternalServerError
	result := struct {
		Message string
	}{
		Message: fmt.Sprintf("Error: %v", err),
	}
	return c.RenderJSON(result)
}
