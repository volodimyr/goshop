package actions

import (
	"bytes"
	"goshop/api/models"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

func AddItemToTheOrder(c buffalo.Context) error {
	db := models.DB
	ord := &models.Order{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("ID must be valid."))
	}
	ord.ID = id
	// if order not found- nothing to update
	err = db.Find(ord, id)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusNotFound, r.String("Nothing to update"))
	}
	// what's the item
	it := &models.Item{}
	if err := c.Bind(it); err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := it.ValidateUpdateOrders(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	//update order
	t := time.Now()
	ord.Updated = t
	ord.Sum += it.Price * it.Count
	err = db.Update(ord, "id", "created")
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(ord))
}

func ExportOrders(c buffalo.Context) error {
	db := models.DB
	ords := &models.Orders{}
	err := db.All(ords)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	w := c.Response()
	w.Header().Set("Content-Disposition", "attachment; filename=orders.json")
	w.Header().Set("Content-Type", w.Header().Get("Content-Type"))
	if _, err := io.Copy(w, bytes.NewReader([]byte(ords.String()))); err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.String("Sent"))
}

func OrdersList(c buffalo.Context) error {
	db := models.DB
	ords := &models.Orders{}
	q := db.PaginateFromParams(c.Params())
	err := q.All(ords)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	c.Set("pagination", q.Paginator)
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
	return c.Render(201, r.JSON(ord))
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
