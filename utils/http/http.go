package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	// Request request parameter
	Request struct {
		Protocol string
		Host     string
		Port     int
		Path     string
		Body     interface{}
	}

	// Response response body
	Response struct {
		Status       bool        `json:"status"`
		Code         int         `json:"code"`
		ErrorMessage string      `json:"errMessage,omitempty"`
		Data         interface{} `json:"data,omitempty"`
		Message      interface{} `json:"message,omitempty"`
	}
)

func Post(req *Request) (Response, error) {
	body, err := json.Marshal(req.Body)
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}

	resp, err := http.Post(fmt.Sprintf("%s://%s:%d%s", req.Protocol, req.Host, req.Port, req.Path), "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}

	defer resp.Body.Close()
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	return response, nil
}

func PostDynamic(req *Request) (interface{}, error) {
	body, err := json.Marshal(req.Body)
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}

	resp, err := http.Post(fmt.Sprintf("%s://%s:%d%s", req.Protocol, req.Host, req.Port, req.Path), "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}

	defer resp.Body.Close()
	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	return response, nil
}

//PostWithHeader call to specific url using POST methods with Headers
func PostWithHeader(req *Request, headers map[string]string) (*Response, error) {
	body, err := json.Marshal(req.Body)
	if err != nil {
		log.Printf("[http]: error message: %v", err)
		return nil, err
	}
	log.Printf("[http]: req: %+v headers: %+v", req, headers)

	client := &http.Client{}

	if req.Protocol == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
		client = &http.Client{Transport: tr}
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s://%s%s", req.Protocol, req.Host, req.Path), bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[http]: error message: %v", err)
		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("[http]: error message: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("[http]: error message: %v", err)
		return nil, err
	}
	return &response, nil
}
