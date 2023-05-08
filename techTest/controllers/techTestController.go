package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"tehcTest/entities"
	"tehcTest/services"
)

type TechTestController struct {
	DB              *mongo.Database
	TechTestService *services.TechTestService
}

type RequestFibonacci struct {
	N int `json:"n"`
}

type RequestCombination struct {
	N uint `json:"n"`
	R uint `json:"r"`
}

type RequestCompany struct {
	CompanyName     string `json:"company_name"`
	TelephoneNumber string `json:"telephone_number"`
	Address         string `json:"address"`
}

type RequestEmployee struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	JobTitle    string `json:"jobtitle"`
}

func output(status int, code string, data interface{}, message string) gin.H {
	return gin.H{
		"status":  status,
		"code":    code,
		"data":    data,
		"message": message,
	}
}

func (tc *TechTestController) Fibonacci(ctx *gin.Context) {
	var requestFibonacci RequestFibonacci
	if err := ctx.ShouldBindJSON(&requestFibonacci); err != nil || requestFibonacci.N == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, "n is required"))
		return
	}

	result := tc.TechTestService.Fibonacci(requestFibonacci.N)
	ctx.JSON(http.StatusOK, output(200, "200", result, "Success"))
}

func (tc *TechTestController) Combination(ctx *gin.Context) {
	var requestCombination RequestCombination
	err := ctx.ShouldBindJSON(&requestCombination)
	if err != nil || requestCombination.N == 0 || requestCombination.R == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, "n or r is required"))
		return
	}

	if requestCombination.R > requestCombination.N {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, "the value of r must be less than n"))
		return
	}

	result := tc.TechTestService.Combination(requestCombination.N, requestCombination.R)
	ctx.JSON(http.StatusOK, output(200, "200", result, "success"))
}

func (tc *TechTestController) GetCountries(ctx *gin.Context) {

	result, err := tc.TechTestService.GetCountries()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, output(500, "500", nil, err.Error()))
	}

	ctx.JSON(http.StatusOK, output(200, "200", result, "success"))
}

func (tc *TechTestController) AddCompany(ctx *gin.Context) {
	var requestCompany RequestCompany
	if err := ctx.ShouldBindJSON(&requestCompany); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}

	company := entities.Company{
		ID:              primitive.NewObjectID(),
		TelephoneNumber: requestCompany.TelephoneNumber,
		CompanyName:     requestCompany.CompanyName,
		Address:         requestCompany.Address,
	}

	id, err := tc.TechTestService.InsertCompany(company)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", company, err.Error()))
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	ctx.JSON(http.StatusCreated, output(201, "201", data, "Success"))
}

func (tc *TechTestController) GetCompenies(ctx *gin.Context) {
	datas, err := tc.TechTestService.FindCompany()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", datas, err.Error()))
		return
	}
	result := make(map[string]interface{})
	result["count"] = len(datas)
	result["rows"] = datas
	ctx.JSON(http.StatusOK, output(200, "200", result, "Success"))
}

func (tc *TechTestController) SetCompanyActive(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := tc.TechTestService.UpdateCompanyActive(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}
	result := make(map[string]interface{})
	result["id"] = data.ID
	result["is_active"] = data.IsActive
	ctx.JSON(http.StatusOK, output(200, "200", result, "Success"))
}

func (tc *TechTestController) AddEmployee(ctx *gin.Context) {
	companyId := ctx.Param("company_id")

	var requestEmployee RequestEmployee
	if err := ctx.ShouldBindJSON(&requestEmployee); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}

	employee := entities.Employee{
		ID:          primitive.NewObjectID(),
		Name:        requestEmployee.Name,
		Email:       requestEmployee.Email,
		PhoneNumber: requestEmployee.PhoneNumber,
		JobTitle:    requestEmployee.JobTitle,
	}

	employeeId, err := tc.TechTestService.InsertEmployee(companyId, employee)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}
	data := make(map[string]interface{})
	data["id"] = employeeId
	data["company_id"] = companyId

	ctx.JSON(http.StatusCreated, output(201, "201", data, "Success"))
}

func (tc *TechTestController) GetEmployeesByCompanyId(ctx *gin.Context) {
	companyId := ctx.Param("id")
	datas, err := tc.TechTestService.FindEmployeesByCompanyId(companyId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", datas, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, output(200, "200", datas, "Success"))
}

func (tc *TechTestController) UpdateEmployee(ctx *gin.Context) {
	companyId := ctx.Param("id")
	employeeId := ctx.Param("employee_id")

	var requestEmployee RequestEmployee
	if err := ctx.ShouldBindJSON(&requestEmployee); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}

	updateEmployee := entities.Employee{
		Name:        requestEmployee.Name,
		Email:       requestEmployee.Email,
		PhoneNumber: requestEmployee.PhoneNumber,
		JobTitle:    requestEmployee.JobTitle,
	}

	err := tc.TechTestService.UpdateEmployee(companyId, employeeId, updateEmployee)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}
	data := make(map[string]interface{})
	data["id"] = employeeId
	data["company_id"] = companyId

	ctx.JSON(http.StatusCreated, output(201, "201", data, "Success"))
}

func (tc *TechTestController) GetEmployeeById(ctx *gin.Context) {
	employeeId := ctx.Param("id")
	data, err := tc.TechTestService.FindEmployeesByEmployeeId(employeeId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, output(200, "200", data, "Success"))
}

func (tc *TechTestController) DeleteEmployee(ctx *gin.Context) {
	employeeId := ctx.Param("id")
	err := tc.TechTestService.DeleteEmployeesByEmployeeId(employeeId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, output(400, "400", nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
