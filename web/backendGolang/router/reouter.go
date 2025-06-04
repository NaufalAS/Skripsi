package router

import (
	"net/http"
	"skripsi/app"
	analiticscontroller "skripsi/controller/analitics"
	datacontroller "skripsi/controller/data"
	usercontroller "skripsi/controller/user"
	"skripsi/helper"
	"skripsi/model"
	query_builder_data "skripsi/query_builder/data_qeury"
	analiticsrepo "skripsi/repository/analitics"
	datarepo "skripsi/repository/data"
	userrepo "skripsi/repository/user"
	analiticsservice "skripsi/service/analitics"
	dataservice "skripsi/service/data"
	userservice "skripsi/service/user"

	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// var tokenUseCase helper.TokenUseCase

func RegisterUserRoute(prefix string, e *echo.Echo) {
	db := app.DBConnection()

	userRepo := userrepo.NewAuthRepository(db)
	userService := userservice.NewSektorUsahaService(userRepo)
	userController := usercontroller.NewSektorUsahaController(userService)

	dataQueryBuilder := query_builder_data.NewDataQueryBuilder(db)
	datarepo := datarepo.NewDataRepository(db, dataQueryBuilder)
	dataService := dataservice.NewSektorDataService(datarepo)
	dataController := datacontroller.NewDataController(dataService)

	analiticsRepo := analiticsrepo.NewAnalticsRepository(db)
	analiticsService := analiticsservice.NewAnaliticsService(analiticsRepo)
	analiticsController := analiticscontroller.NewAnaliticsController(analiticsService)

	
	g := e.Group(prefix)
	authRoute := g.Group("/user")
	authRoute.GET("/list", userController.GetListUserController)
	authRoute.POST("/post", userController.PostUserController)
	authRoute.POST("/login", userController.LoginUserController)
	authRoute.GET("/:id", userController.GetUserByIdController)
	authRoute.PUT("/update/:id", userController.UpdateUserByIdController, JWTProtection())
	authRoute.DELETE("/delete/:id", userController.DeleteProdukId)
	authRoute.PUT("/password/:id", userController.UpdatePasswordController, JWTProtection())

	dataRoute := g.Group("/data")
	dataRoute.POST("/post", dataController.PostDataController)
	dataRoute.GET("/list", dataController.GetListDataController, JWTProtection())
	dataRoute.GET("/:id", dataController.GetDataByIdController, JWTProtection())
	dataRoute.DELETE("/delete/:id", dataController.DeleteDataId, JWTProtection())
	dataRoute.PUT("/update/:id", dataController.UpdateDataByIdController, JWTProtection())

	analiticsRoute := g.Group("/analitics")
	analiticsRoute.GET("/list", analiticsController.GetAnalyticsDataController)
}
func JWTProtection() echo.MiddlewareFunc {
		return echojwt.WithConfig(echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(helper.JwtCustomClaims)
			},
			SigningKey: []byte(os.Getenv("SECRET_KEY")),
			ErrorHandler: func(c echo.Context, err error) error {
				return c.JSON(http.StatusUnauthorized, model.ResponseToClient(http.StatusUnauthorized, "unauthorized", nil))
			},
		})
	}