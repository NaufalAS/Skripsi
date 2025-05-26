package router

import (
	"net/http"
	"skripsi/app"
	datacontroller "skripsi/controller/data"
	usercontroller "skripsi/controller/user"
	"skripsi/helper"
	"skripsi/model"
	datarepo "skripsi/repository/data"
	userrepo "skripsi/repository/user"
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

	datarepo := datarepo.NewDataRepository(db)
	dataService := dataservice.NewSektorDataService(datarepo)
	dataController := datacontroller.NewDataController(dataService)

	
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
	dataRoute.POST("/post", dataController.PostDataController, JWTProtection())
	dataRoute.GET("/list", dataController.GetListDataController, JWTProtection())
	dataRoute.GET("/:id", dataController.GetDataByIdController, JWTProtection())
	dataRoute.DELETE("/delete/:id", dataController.DeleteDataId, JWTProtection())
	dataRoute.PUT("/update/:id", dataController.UpdateDataByIdController, JWTProtection())
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