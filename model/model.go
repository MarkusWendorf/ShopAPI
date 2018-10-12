package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type Product struct {
	Id         string            `json:"id" bson:"id,omitempty"`
	Category   string            `json:"category,omitempty"`
	Pname      string            `json:"pname,omitempty"`
	Price      int               `json:"price,omitempty"`
	Quantity   int               `json:"quantity"`
	ImgUrl     string            `json:"imgurl,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (p Product) String() string {
	return p.Pname + ", Price: " + strconv.Itoa(p.Price)
}

type AutocompleteProduct struct {
	Id        string `json:"id"`
	Category  string `json:"category"`
	Pname     string `json:"pname"`
	Highlight string `json:"highlight"`
}

type User struct {
	Email    string        `json:"email"`
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Password string        `json:"password"`
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type Cart struct {
	Items    []CartItem    `json:"items"`
	Total    int           `json:"total"`
	Shipping int           `json:"shipping"`
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Expires  int64         `json:"expires,omitempty"`
	UserId   bson.ObjectId `json:"user_id"`
}

type Order struct {
	Id      bson.ObjectId      `json:"id" bson:"_id,omitempty"`
	Cart    Cart               `json:"cart"`
	Address AddressInformation `json:"address"`
}

type Address struct {
	City     string `json:"city"`
	CityCode string `json:"citycode"`
	Street   string `json:"street"`
	Name     string `json:"name"`
}

type AddressInformation struct {
	Shipping Address `json:"shipping"`
	Billing  Address `json:"billing"`
}

type Response struct {
	Data  interface{}            `json:"data"`
	Error string                 `json:"error,omitempty"`
	Links map[string]interface{} `json:"links"`
}

type JwtClaims struct {
	UserId bson.ObjectId `json:"user_id"`
	Email  string        `json:"email"`
	Exp    int64         `json:"exp"`
}

func (j JwtClaims) Valid() error {

	if time.Now().Unix() > j.Exp {
		return fmt.Errorf("token has expired (was valid until: %v)", j.Exp)
	}

	return nil
}
