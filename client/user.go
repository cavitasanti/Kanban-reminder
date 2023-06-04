package client

import (
	"a21hc3NpZ25tZW50/config"
	"a21hc3NpZ25tZW50/entity"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type UserClient interface {
	Login(email, password string) (userId int, respCode int, err error)
	Register(fullname, email, password string) (userId int, respCode int, err error)
	GetUser(userId string) (entity.User, error)
	UpdateUser(userId, fullname string) (respCode int, err error)
	DeleteUser(userId string) (respCode int, err error)
}

type userClient struct {
}

func NewUserClient() *userClient {
	return &userClient{}
}

func (u *userClient) Login(email, password string) (userId int, respCode int, err error) {
	datajson := map[string]string{
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return 0, -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/users/login"), bytes.NewBuffer(data))
	if err != nil {
		return 0, -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, -1, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}

	if err != nil {
		return 0, -1, err
	} else {
		json.Unmarshal(b, &result)

		if result["user_id"] != nil {
			return int(result["user_id"].(float64)), resp.StatusCode, nil
		} else {
			return 0, resp.StatusCode, errors.New(string(result["error"].(string)))
		}
	}
}

func (u *userClient) Register(fullname, email, password string) (userId int, respCode int, err error) {
	datajson := map[string]string{
		"fullname": fullname,
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return 0, -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/users/register"), bytes.NewBuffer(data))
	if err != nil {
		return 0, -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, -1, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}

	if err != nil {
		return 0, -1, err
	} else {
		json.Unmarshal(b, &result)

		if result["user_id"] != nil {
			return int(result["user_id"].(float64)), resp.StatusCode, nil
		} else {
			return 0, resp.StatusCode, errors.New(string(result["error"].(string)))
		}
	}
}

func (u *userClient) GetUser(userId string) (entity.User, error) {
	client, err := GetClientWithCookie(userId)
	if err != nil {
		return entity.User{}, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/users/get"), nil)
	if err != nil {
		return entity.User{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return entity.User{}, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	var result entity.User
	if err != nil {
		return entity.User{}, err
	} else {
		json.Unmarshal(b, &result)

		return result, nil
	}
}

func (u *userClient) UpdateUser(userId, fullname string) (respCode int, err error) {
	datajson := map[string]string{
		"fullname": fullname,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	client, err := GetClientWithCookie(userId)
	if err != nil {
		return -1, err
	}
	log.Print(data)
	req, err := http.NewRequest("PUT", config.SetUrl("/api/v1/users/update"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (u *userClient) DeleteUser(userId string) (respCode int, err error) {
	req, err := http.NewRequest("DELETE", config.SetUrl("/api/v1/users/delete"+userId), nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
