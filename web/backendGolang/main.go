// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"skripsi/router"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/joho/godotenv"
// 	"github.com/labstack/echo/v4"
// )

// // CustomValidator untuk Echo
// type CustomValidator struct {
//     validator *validator.Validate
// }

// func (cv *CustomValidator) Validate(i interface{}) error {
//     return cv.validator.Struct(i)
// }

// func main() {
//     // Load environment variables
//     if err := godotenv.Load(".env"); err != nil {
//         log.Fatal("Error loading .env file!")
//     }

//     // Ambil port dari environment variable
//     port := os.Getenv("PORT")
//     if port == "" {
//         port = "8001" // Default port jika tidak ada di .env
//     }

//     // Inisialisasi Echo
//     e := echo.New()

//     // Set custom validator
//     e.Validator = &CustomValidator{validator: validator.New()}

//     // Daftarkan route
//     router.RegisterUserRoute("/api", e)

//     // Cek koneksi berjalan
//     fmt.Println("Server running on port:", port)	

//     // Start server
//     if err := e.Start(":" + port); err != nil {
//         log.Fatal("Error starting server:", err)
//     }
// }


package main

import (
	"fmt"
	"log"
	"os"
	"skripsi/router"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator untuk Echo
type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    // Load environment variables
    if err := godotenv.Load(".env"); err != nil {
        log.Fatal("Error loading .env file!")
    }

    // Ambil port dari environment variable
    port := os.Getenv("PORT")
    if port == "" {
        port = "8001" // Default port jika tidak ada di .env
    }

    // Inisialisasi Echo
    e := echo.New()

    // Set custom validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Enable CORS
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend localhost:5173
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    }))

    // Daftarkan route
    router.RegisterUserRoute("/api", e)
    e.Static("/profile", "public/profile")
    e.Static("/pelanggaran", "public/pelanggaran")


    // Cek koneksi berjalan
    fmt.Println("Server running on port:", port)	

    // Start server
    if err := e.Start(":" + port); err != nil {
        log.Fatal("Error starting server:", err)
    }
}
