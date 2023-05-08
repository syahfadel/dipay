package routers

import (
	"context"
	"tehcTest/controllers"
	"tehcTest/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartService(db *mongo.Database, ctx *context.Context) *gin.Engine {

	techService := services.TechTestService{
		DB: db,
	}

	techController := controllers.TechTestController{
		DB:              db,
		TechTestService: &techService,
	}

	app := gin.Default()
	app.POST("/api/fibonacci", techController.Fibonacci)
	app.POST("/api/combination", techController.Combination)
	app.GET("/api/countries", techController.GetCountries)

	companies := app.Group("/api/companies")
	{
		companies.POST("/", techController.AddCompany)
		companies.GET("/", techController.GetCompenies)
		companies.PUT("/:id/set_active", techController.SetCompanyActive)
		companies.POST("/:company_id/employees", techController.AddEmployee)
		companies.GET("/:id/employees", techController.GetEmployeesByCompanyId)
		companies.PUT("/:id/employees/:employee_id", techController.UpdateEmployee)
	}

	employees := app.Group("/api/employees")
	{
		employees.GET("/:id", techController.GetEmployeeById)
		employees.DELETE("/:id", techController.DeleteEmployee)
	}

	return app
}
