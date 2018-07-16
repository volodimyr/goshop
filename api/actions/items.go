package actions

import (
	"goshop/api/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

// ItemsList default implementation.
func ItemsList(c buffalo.Context) error {
	db := models.DB
	its := &models.Items{}
	err := db.All(its)

	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(its))
}

// ItemsIndex default implementation.
func ItemsIndex(c buffalo.Context) error {
	db := models.DB
	it := &models.Item{}
	uID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Bad Request"))
	}
	err = db.Find(it, uID)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusNotFound, r.String("No item has been found"))
	}
	return c.Render(200, r.JSON(it))
}

// ItemsDelete remove category by id
func ItemsDelete(c buffalo.Context) error {
	db := models.DB
	it := &models.Item{}
	uID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Bad Request"))
	}
	it.ID = uID
	err = db.Destroy(it)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(it))
}

func ItemsCreate(c buffalo.Context) error {
	db := models.DB
	it := &models.Item{}
	if err := c.Bind(it); err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := it.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	uID, err := uuid.NewV4()
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	it.ID = uID
	err = db.Create(it)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(it))
}

func ItemsUpdate(c buffalo.Context) error {
	db := models.DB
	it := &models.Item{}
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("ID must be valid."))
	}
	it.ID = id
	if err := c.Bind(it); err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := it.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	//exclude id
	err = db.Update(it, "id")
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(it))
}
