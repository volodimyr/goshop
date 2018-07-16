package actions

import (
	"goshop/api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

func OrdersList(c buffalo.Context) error {
	db := models.DB
	ords := &models.Orders{}
	err := db.All(ords)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(ords))
}

func OrdersIndex(c buffalo.Context) error {
	db := models.DB
	ord := &models.Order{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("ID must be valid"))
	}
	err = db.Find(ord, id)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusNotFound, r.String("No order has been found"))
	}
	return c.Render(200, r.JSON(ord))
}

// OrdersCreate default implementation.
func OrdersCreate(c buffalo.Context) error {
	db := models.DB
	ord := &models.Order{}
	if err := c.Bind(ord); err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := ord.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	t := time.Now()
	ord.Created, ord.Updated = t, t
	err := db.Create(ord, "id")
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(ord))
}

// OrdersUpdate default implementation.
func OrdersUpdate(c buffalo.Context) error {
	db := models.DB
	ord := &models.Order{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("ID must be valid."))
	}
	ord.ID = id
	if err := c.Bind(ord); err != nil {
		return errors.WithStack(err)
		// return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := ord.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	t := time.Now()
	ord.Updated = t
	err = db.Update(ord, "id", "created")
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(ord))
}
