package actions

import (
	"encoding/json"
	"goshop/api/models"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

func ImportCategories(c buffalo.Context) error {
	f, err := c.File("categories")
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't upload file."))
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't read file."))
	}

	cgs := &models.Categories{}
	err = json.Unmarshal(buf, cgs)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.String("Couldn't unmarshal json"))
	}

	db := models.DB
	for _, cg := range *cgs {
		cg.ID, err = uuid.NewV4()
		if err != nil {
			// return errors.WithStack(err)
			return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
		}
	}

	err = db.Create(cgs)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(cgs))
}

// CategoriesList default implementation.
func CategoriesList(c buffalo.Context) error {
	db := models.DB
	cgs := &models.Categories{}
	q := db.PaginateFromParams(c.Params())
	err := q.All(cgs)

	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.JSON(cgs))
}

// CategoriesIndex default implementation.
func CategoriesIndex(c buffalo.Context) error {
	db := models.DB
	cg := &models.Category{}
	uID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Bad Request"))
	}
	err = db.Find(cg, uID)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusNotFound, r.String("No category has been found"))
	}
	return c.Render(200, r.JSON(cg))
}

// CategoriesDelete remove category by id
func CategoriesDelete(c buffalo.Context) error {
	db := models.DB
	cg := &models.Category{}
	uID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Bad Request"))
	}
	cg.ID = uID
	err = db.Destroy(cg)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(cg))
}

// CategoriesCreate create new Category
func CategoriesCreate(c buffalo.Context) error {
	db := models.DB
	cg := &models.Category{}
	if err := c.Bind(cg); err != nil {
		// return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := cg.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	uID, err := uuid.NewV4()
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	cg.ID = uID
	err = db.Create(cg)
	if err != nil {
		// 		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(201, r.JSON(cg))
}

// CategoriesUpdate update the existing or create new one if doesn't exist
func CategoriesUpdate(c buffalo.Context) error {
	db := models.DB
	cg := &models.Category{}
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("ID must be valid."))
	}
	cg.ID = id
	if err := c.Bind(cg); err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String("Couldn't bind payload."))
	}
	vErr, _ := cg.Validate(db)
	if vErr.HasAny() != false {
		//		return errors.WithStack(err)
		return c.Render(http.StatusBadRequest, r.String(vErr.String()))
	}
	err = db.Update(cg, "id")
	if err != nil {
		//		return errors.WithStack(err)
		return c.Render(http.StatusInternalServerError, r.String("Internal server error"))
	}
	return c.Render(200, r.JSON(cg))
}
