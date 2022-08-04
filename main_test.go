package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/yan.ren/go-rest-api-mysql/model"
)

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest(echo.GET, "http://localhost:8080/users", strings.NewReader(""))
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	client := http.Client{}
	response, err := client.Do(req)
	assertNoError(t, err)
	assertEqual(t, http.StatusOK, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	assertNoError(t, err)
	assertEqual(t, `{"users":[{"id":1,"name":"Solomon"},{"id":2,"name":"Menelik"}]}`, strings.Trim(string(byteBody), "\n"))

	response.Body.Close()
}

func TestGetUserById(t *testing.T) {
	req, err := http.NewRequest(echo.GET, "http://localhost:8080/users/1", strings.NewReader(""))
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	client := http.Client{}
	response, err := client.Do(req)
	assertNoError(t, err)
	assertEqual(t, http.StatusOK, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	assertNoError(t, err)
	assertEqual(t, `{"user":{"id":1,"name":"Solomon"}}`, strings.Trim(string(byteBody), "\n"))

	response.Body.Close()
}

type UserResponse struct {
	User model.User `json:"user"`
}

func TestPOST(t *testing.T) {
	reqStr := `{
		"name": "test"
	  }`
	req, err := http.NewRequest(echo.POST, "http://localhost:8080/user", strings.NewReader(reqStr))
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	assertNoError(t, err)
	assertEqual(t, http.StatusCreated, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	assertNoError(t, err)

	var res UserResponse
	if error := json.Unmarshal(byteBody, &res); err != nil {
		panic(error)
	}
	if res.User.ID != 3 {
		t.Logf("new user id is not matched, received: %d", res.User.ID)
		t.Fail()
	}
	if res.User.Name != "test" {
		t.Logf("new user name is not matched, received: %s", res.User.Name)
		t.Fail()
	}

	response.Body.Close()
}

func TestPATCH(t *testing.T) {
	reqStr := `{
		"name": "test"
	  }`
	req, err := http.NewRequest(echo.PATCH, "http://localhost:8080/users/1", strings.NewReader(reqStr))
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	assertNoError(t, err)
	assertEqual(t, http.StatusCreated, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	assertNoError(t, err)

	var res UserResponse
	if error := json.Unmarshal(byteBody, &res); err != nil {
		panic(error)
	}
	if res.User.ID != 1 {
		t.Logf("patched user id is not matched, received: %d", res.User.ID)
		t.Fail()
	}
	if res.User.Name != "test" {
		t.Logf("patched user name is not matched, received: %s", res.User.Name)
		t.Fail()
	}

	response.Body.Close()
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
}
