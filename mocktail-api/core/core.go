package core

import (
	"mocktail-api/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type Api struct {
	ID       uint           `gorm:"primary_key;auto_increment;not_null"`
	Endpoint string         `validate:"required"`
	Method   string         `validate:"is-method-allowed"`
	Key      string         `gorm:"unique;not null"`
	Response datatypes.JSON `validate:"required"`
}

type apiResponse struct {
	Message string `json:""`
}

type Apis struct {
	Apis []Api `validate:"required"`
}

func GetApis(c *fiber.Ctx) error {
	db := database.DBConn
	var apis []Api
	db.Find(&apis)
	return c.JSON(apis)
}

func CreateApi(c *fiber.Ctx) error {
	api := new(Api)
	if err := c.BodyParser(api); err != nil {
		res := apiResponse{err.Error()}
		return c.Status(503).JSON(res)
	}
	if err := InsertApi(api); err != nil {
		res := apiResponse{err.Error()}
		return c.Status(400).JSON(res)
	}
	return c.JSON(api)
}

func InsertApi(api *Api) error {
	db := database.DBConn
	api.Key = api.Method + api.Endpoint
	validate := validator.New()
	validate.RegisterValidation("is-method-allowed", isApiHTTPMethodValid)
	if err := validate.Struct(api); err != nil {
		return err
	}
	if err := db.Create(&api).Error; err != nil {
		return err
	}

	return nil
}

func DeleteApiByKey(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var api Api
	res := apiResponse{"completed"}
	err := db.Unscoped().Delete(&api, "ID = ? ", id).Error
	if err != nil {
		res := apiResponse{err.Error()}
		return c.Status(400).JSON(res)
	}
	return c.JSON(res)
}

func ExportApis(c *fiber.Ctx) error {
	db := database.DBConn
	var apis []Api
	db.Find(&apis)
	return c.JSON(apis)
}
func ImportApis(c *fiber.Ctx) error {

	apis := new(Apis)
	res := apiResponse{"completed"}
	if err := c.BodyParser(apis); err != nil {
		res := apiResponse{err.Error()}
		return c.Status(400).JSON(res)
	}

	for i := 0; i < len(apis.Apis); i++ {
		InsertApi(&apis.Apis[i])
	}
	return c.Status(200).JSON(res)
}

func isApiHTTPMethodValid(fl validator.FieldLevel) bool {
	HTTPMethodList := [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, b := range HTTPMethodList {
		if b == fl.Field().String() {
			return true
		}
	}
	return false
}
