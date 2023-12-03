package mocktail

import (
	"mocktail-api/database"

	"github.com/kataras/iris/v12"
	"gorm.io/datatypes"
)

type Api struct {
	ID       uint           `gorm:"primary_key;auto_increment;not_null"`
	BaseUrl  string         `gorm:"not_null"`
	Endpoint string         `validate:"required"`
	Method   string         `validate:"is-method-allowed"`
	Key      string         `gorm:"unique;not null"`
	Qstring  string         `gorm:""`
	Response datatypes.JSON `validate:"required"`
}

func MockApiHandler(ctx iris.Context) {
	db := database.DBConn
	var api Api
	// key := strings.Replace(string(ctx.Path()), "/mocktail/", ctx.Method(), 1)
	key := ctx.Method() + ctx.Request().RequestURI
	db.Where("key = ?", key).First(&api)
	if api.Key == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "API not found..."})
	}
	ctx.JSON(api.Response)
}
