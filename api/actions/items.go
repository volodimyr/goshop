package actions

import "github.com/gobuffalo/buffalo"

// ItemsList default implementation.
func ItemsList(c buffalo.Context) error {
	return c.Render(200, r.JSON("items/list.json"))
}

// ItemsIndex default implementation.
func ItemsIndex(c buffalo.Context) error {
	return c.Render(200, r.JSON("items/index.json"))
}
