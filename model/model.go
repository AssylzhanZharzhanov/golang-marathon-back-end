package model

import (
	"time"
)

type Token struct {
	UserID int    `json:"userId"`
	Token  string `json:"token"`
}

type User struct {
	FirstName        string    `json:"firstName"`
	SecondName       string    `json:"secondName"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	RegistrationTime time.Time `json:"registrationTime,omitempty"`
}

type Product struct {
	UserID int `json:"userId"`
}

type Search struct {
	Parameter string `json:"parameter"`
}

type Marathon struct {
	MarathonID   int    `json:"marathonId"`
	MarathonName string `json:"marathonName"`
}

type ImageWatermark struct {
	ImageID    int    `json:"imageId"`
	Image      string `json:"image"`
	MarathonID int    `json:"marathonId"`
}

type Image struct {
	ImageID    int    `json:"imageId"`
	Image      string `json:"image"`
	MarathonID int    `json:"marathonId"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     string `json:"result"`
}

type ResponseSignup struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     User   `json:"result"`
}

type ResponseToken struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     Token  `json:"result"`
}

type ResponseMarathons struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Result     []Marathon `json:"result"`
}

type ResponseImages struct {
	StatusCode int     `json:"statusCode"`
	Message    string  `json:"message"`
	Result     []Image `json:"result"`
}

type ResponseImageWatermark struct {
	StatusCode int              `json:"statusCode"`
	Message    string           `json:"message"`
	Result     []ImageWatermark `json:"result"`
}

type ResponseSuccess struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     string `json:"result"`
}

// type MongoImage struct {
// 	ID          bson.ObjectId `json:"_id"`
// 	Author      string        `json:"author"`
// 	Caption     string        `json:"caption"`
// 	ContentType string        `json:"contentType"`
// 	DateTime    string        `json:"dateTime"`
// 	FileID      bson.ObjectId `json:"fileID"`
// 	FileSize    int64         `json:"fileSize"`
// 	Name        string        `json:"name"`
// }
