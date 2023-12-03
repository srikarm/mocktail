package core

import (
	"mocktail-api/database"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"gorm.io/datatypes"
)

type Api struct {
	ID       uint           `gorm:"primary_key;auto_increment;not_null"`
	BaseUrl  string         `gorm:"not_null"`
	Endpoint string         `validate:"required"`
	Method   string         `validate:"is-method-allowed"`
	Key      string         `gorm:"unique;not null"`
	Qstring  string			`gorm:""`
	Response datatypes.JSON `validate:"required"`
}

type Apis struct {
	Apis []Api `validate:"required"`
}

func GetApis(ctx iris.Context) {
	db := database.DBConn
	var apis []Api
	db.Find(&apis)
	ctx.JSON(apis)
}

func CreateApi(ctx iris.Context) {
	api := new(Api)
	if err := ctx.ReadBody(api); err != nil {
		ctx.StopWithJSON(iris.StatusServiceUnavailable, iris.Map{"response": err.Error()})
		return
	}
	if err := InsertApi(api); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"response": err.Error()})
		return
	}
	ctx.JSON(iris.Map{"response": "API creation has been successfull", "message": "completed"})
}

func InsertApi(api *Api) error {
	db := database.DBConn
	params, _ := url.Parse(api.Endpoint)
	key := []string{api.Method,api.BaseUrl,api.Endpoint}
	api.Key = strings.Join(key,"/")
	api.Qstring = params.RawQuery
//	api.Endpoint = params.Path
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

func DeleteApiByKey(ctx iris.Context) {
	id := ctx.Params().Get("id")
	db := database.DBConn
	var api Api
	err := db.Unscoped().Delete(&api, "ID = ? ", id).Error
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"response": err.Error()})
		return
	}
	ctx.JSON(iris.Map{"response": "API has been removed", "message": "completed"})
}

func ExportApis(ctx iris.Context) {
	db := database.DBConn
	var apis []Api
	db.Find(&apis)
	ctx.JSON(apis)
}
func ImportApis(ctx iris.Context) {

	apis := new(Apis)
	if err := ctx.ReadBody(apis); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{"response": err.Error()})
		return
	}

	for i := 0; i < len(apis.Apis); i++ {
		InsertApi(&apis.Apis[i])
	}
	ctx.JSON(iris.Map{"response": "APIs have been imported", "message": "completed"})
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
