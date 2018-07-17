package grifts

import (
	"fmt"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		return nil
	})

})

var _ = grift.Add("hello", func(c *grift.Context) error {
	fmt.Println("Hello World!")
	return nil
})
