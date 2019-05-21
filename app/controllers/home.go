package controllers

import (
	"github.com/revel/revel"
)

type Home struct {
	Base
}

func (c Home) Index() revel.Result {
	return c.Render()
}
