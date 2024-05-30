package main

import (
	"HashingTableStudent/model"
	"HashingTableStudent/utils"
	"HashingTableStudent/views"
	"github.com/labstack/echo/v4"
	"strconv"
)

var (
	errorData = ""
)

func main() {
	e := echo.New()
	e.File("/file", "students.txt")
	e.GET("/", func(c echo.Context) error {
		data, err := utils.LoadAllData()
		if err != nil {
			panic(err)
		}
		return views.Index(data, errorData).Render(c.Request().Context(), c.Response())
	})

	stu := e.Group("/student")
	stu.GET("/add", func(c echo.Context) error {
		return views.RowAdd().Render(c.Request().Context(), c.Response())
	})
	stu.GET("/find", func(c echo.Context) error {
		studentNumber := c.FormValue("StudentNumber")
		st, line, err := utils.SearchStudent(studentNumber)
		if err != nil {
			errorData = err.Error()
		} else if line == 0 {
			errorData = "student not found"
			return c.Redirect(302, "/")
		}
		return views.RowShow(*st, errorData).Render(c.Request().Context(), c.Response())
	})
	stu.GET("/edit/:studentNumber", func(c echo.Context) error {
		studentNumber := c.Param("studentNumber")
		st, line, err := utils.SearchStudent(studentNumber)
		if err != nil {
			errorData = err.Error()
			return c.Redirect(302, "/")
		} else if line == 0 {
			errorData = "student not found"
			return c.Redirect(302, "/")
		}
		return views.RowEdit(*st).Render(c.Request().Context(), c.Response())
	})
	stu.POST("/add", func(c echo.Context) error {
		enteringYear, err := strconv.Atoi(c.FormValue("EnteringYear"))
		if err != nil {
			errorData = err.Error()
			return c.Redirect(302, "/")
		}
		gpa, err := strconv.ParseFloat(c.FormValue("GPA"), 64)
		if err != nil {
			errorData = err.Error()
			return c.Redirect(302, "/")
		}
		var std model.Student = model.Student{
			StudentNumber: c.FormValue("StudentNumber"),
			NationalCode:  c.FormValue("NationalCode"),
			Name:          c.FormValue("Name"),
			LastName:      c.FormValue("LastName"),
			EnteringYear:  enteringYear,
			GPA:           gpa,
		}
		utils.AddOrEditStudent(std)
		return c.Redirect(302, "/")
	})

	stu.GET("/delete/:studentNumber", func(c echo.Context) error {
		studentNumber := c.Param("studentNumber")
		_, err := utils.DeleteStudent(studentNumber)
		if err != nil {
			errorData = err.Error()
			return c.Redirect(302, "/")
		}
		return c.Redirect(302, "/")
	})

	stu.GET("/hash", func(c echo.Context) error {
		studentNumber := c.FormValue("hashinput")
		return c.String(200, strconv.Itoa(utils.Hash(studentNumber)))
	})

	e.Logger.Fatal(e.Start(":8001"))
}
