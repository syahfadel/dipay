package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"tehcTest/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TechTestService struct {
	DB *mongo.Database
}

var ctx = context.Background()

func (ts *TechTestService) Fibonacci(n int) string {
	result := "0"
	prev := 0
	current := 1
	for current <= n {
		result += fmt.Sprintf(" %v", current)
		next := current + prev
		prev = current
		current = next
	}

	return result
}

func (ts *TechTestService) Combination(n, r uint) *big.Int {
	return new(big.Int).Div(factorial(n), new(big.Int).Mul(factorial(r), factorial(n-r)))
}

func factorial(n uint) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}

	return new(big.Int).Mul(big.NewInt(int64(n)), factorial(n-1))
}

func (ts *TechTestService) GetCountries() ([]interface{}, error) {
	var country map[string]interface{}
	var countries []interface{}

	res, err := http.Get("https://gist.githubusercontent.com/herysepty/ba286b815417363bfbcc472a5197edd0/raw/aed8ce8f5154208f9fe7f7b04195e05de5f81fda/coutries.json")
	result, err := ioutil.ReadAll(res.Body)
	var datas []map[string]interface{}
	err = json.Unmarshal(result, &datas)

	if err != nil {
		return nil, err
	}

	for _, d := range datas {
		country = make(map[string]interface{})
		country["name"] = d["name"]
		country["region"] = d["region"]
		country["timezones"] = d["timezones"]

		countries = append(countries, country)
	}

	return countries, nil
}

func (ts *TechTestService) InsertCompany(company entities.Company) (interface{}, error) {
	collection := ts.DB.Collection("companies")
	if err := company.BeforeCreate(); err != nil {
		return nil, err
	}
	res, err := collection.InsertOne(ctx, company)
	if err != nil {
		return nil, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (ts *TechTestService) FindCompany() ([]entities.Company, error) {
	var result []entities.Company
	collection := ts.DB.Collection("companies")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []entities.Company{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []entities.Company{}, err
	}

	return result, nil
}

func (ts *TechTestService) UpdateCompanyActive(idHex string) (entities.Company, error) {
	var result entities.Company
	collection := ts.DB.Collection("companies")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return entities.Company{}, nil
	}

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"is_active": true}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return entities.Company{}, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return entities.Company{}, errors.New(fmt.Sprintf("id %v not exist", idHex))
	}

	return result, nil

}

func (ts *TechTestService) InsertEmployee(companyIdHex string, employee entities.Employee) (interface{}, error) {
	collectionCompanies := ts.DB.Collection("companies")
	collectionEmployee := ts.DB.Collection("employees")

	// Check company id exist or not
	var company entities.Company
	companyId, err := primitive.ObjectIDFromHex(companyIdHex)
	if err != nil {
		return entities.Employee{}, nil
	}

	filterCompany := bson.M{"_id": companyId}

	err = collectionCompanies.FindOne(ctx, filterCompany).Decode(&company)
	if err == mongo.ErrNoDocuments {
		return entities.Company{}, errors.New(fmt.Sprintf("id %v not exist", companyIdHex))
	}

	employee.CompanyID = companyId

	if err := employee.BeforeCreate(); err != nil {
		return nil, err
	}

	res, err := collectionEmployee.InsertOne(ctx, employee)
	if err != nil {
		return nil, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (ts *TechTestService) FindEmployeesByCompanyId(idHex string) ([]entities.Employee, error) {
	var results []entities.Employee
	companyId, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"company_id": companyId}
	collection := ts.DB.Collection("employees")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (ts *TechTestService) UpdateEmployee(companyIdHex, employeeIdHex string, updateEmployee entities.Employee) error {
	companyId, err := primitive.ObjectIDFromHex(companyIdHex)
	if err != nil {
		return err
	}
	employeeId, err := primitive.ObjectIDFromHex(employeeIdHex)
	if err != nil {
		return err
	}

	updateEmployee.CompanyID = companyId

	filter := bson.M{
		"_id":        employeeId,
		"company_id": companyId,
	}

	collection := ts.DB.Collection("employees")

	var currentData entities.Employee
	err = collection.FindOne(ctx, filter).Decode(&currentData)
	if err == mongo.ErrNoDocuments {
		return errors.New(fmt.Sprintf("company id %v or employee id %v not exist", companyId, employeeId))
	}

	if updateEmployee.Name == "" {
		updateEmployee.Name = currentData.Name
	}
	if updateEmployee.Email == "" {
		updateEmployee.Email = currentData.Email
	}
	if updateEmployee.PhoneNumber == "" {
		updateEmployee.PhoneNumber = currentData.PhoneNumber
	}
	if updateEmployee.JobTitle == "" {
		updateEmployee.JobTitle = currentData.JobTitle
	}

	update := bson.M{"$set": bson.M{
		"name":         updateEmployee.Name,
		"email":        updateEmployee.Email,
		"phone_number": updateEmployee.PhoneNumber,
		"jobtitle":     updateEmployee.JobTitle,
	}}

	if err := updateEmployee.BeforeCreate(); err != nil {
		return err
	}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (ts *TechTestService) FindEmployeesByEmployeeId(idHex string) (entities.Employee, error) {
	var result entities.Employee
	employeeId, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return entities.Employee{}, err
	}

	filter := bson.M{"_id": employeeId}
	collection := ts.DB.Collection("employees")
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return entities.Employee{}, errors.New(fmt.Sprintf("employee id %v not exist", employeeId))
	}

	return result, nil
}

func (ts *TechTestService) DeleteEmployeesByEmployeeId(idHex string) error {
	employeeId, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": employeeId}
	collection := ts.DB.Collection("employees")
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
