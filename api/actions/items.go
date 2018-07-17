package actions

import (
	"encoding/json"
	"goshop/api/models"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

func ImportItems(c buffalo.Context) error {
	f, err := c.File("categories")
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't upload file."))
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't read file."))
	}

	its := &models.Items{}
	err = json.Unmarshal(buf, its)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't unmarshal json"))
	}

	db := models.DB
	for _, it := range *its {
		it.ID, err = uuid.NewV4()
		if err != nil {
			// return errors.WithStack(err)
			return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
		}
	}
	err = db.Create(its)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(its))
}

// ItemsList default implementation.
func ItemsList(c buffalo.Context) error {
	db := models.DB
	its := &models.Items{}
	q := db.PaginateFromParams(c.Params())
	err := q.All(its)

	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	c.Set("pagination", q.Paginator)
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
	return c.Render(201, r.JSON(it))
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
