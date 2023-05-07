package controllers

import (
	"agree-market/database"
	"agree-market/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ulule/deepcopier"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var user entity.User
	json.Unmarshal(requestBody, &user)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	database.Connector.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []entity.User
	database.Connector.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var user entity.User
	json.Unmarshal(requestBody, &user)
	if len(user.Email) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(user)
	var fetchedUser entity.User
	database.Connector.Where("email = ?", user.Email).First(&fetchedUser)
	// database.Connector.Find(&fetchedUser)
	fmt.Println(fetchedUser)
	if fetchedUser.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(fetchedUser.Password)
	fmt.Println(user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password))
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fetchedUser.ID,
		"email":   fetchedUser.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []entity.Product
	database.Connector.Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchKeyword := r.URL.Query().Get("search")
	var products []entity.Product

	database.Connector.Where("name LIKE ?", "%"+searchKeyword+"%").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productId := params["id"]

	fmt.Println(productId)
	var product entity.Product
	database.Connector.First(&product, productId)

	if product.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func GetShoppingCart(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var user entity.User
	json.Unmarshal(requestBody, &user)
	if len(user.Email) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(user)
	var fetchedUser entity.User
	database.Connector.Where("email = ?", user.Email).First(&fetchedUser)
	// database.Connector.Find(&fetchedUser)
	fmt.Println(fetchedUser)
	if fetchedUser.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(fetchedUser.Password)
	fmt.Println(user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password))
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Sudah otentikasi
	var userShoppingCart entity.Shopping_Cart
	var product entity.Product
	database.Connector.Where("user_id = ?", fetchedUser.ID).First(&userShoppingCart)
	database.Connector.Where("id = ?", userShoppingCart.ProductID).First(&product)
	deepcopier.Copy(&fetchedUser).To(&userShoppingCart.User)
	deepcopier.Copy(&product).To(&userShoppingCart.Product)
	fmt.Println(product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userShoppingCart)
}
