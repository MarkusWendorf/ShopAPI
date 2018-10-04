package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"shopApi/argon"
	"shopApi/model"
	"time"
)

const pageSize = 24

type Database struct {
	products *mgo.Collection
	users    *mgo.Collection
	carts    *mgo.Collection
	orders   *mgo.Collection
}

func New(products *mgo.Collection, users *mgo.Collection, carts *mgo.Collection, orders *mgo.Collection) Database {
	return Database{products: products, users: users, carts: carts, orders: orders}
}

// ===== Products =====

func (db *Database) GetCategoryNames() []string {

	categories := make([]string, 0)
	db.products.Find(nil).Distinct("category", &categories)

	return categories
}

func (db *Database) GetProduct(id string) (*model.Product, error) {

	query := db.products.Find(bson.M{"id": id})

	var product model.Product
	err := query.One(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// ===== Users =====

func (db *Database) GetUser(email string) (model.User, error) {

	query := db.users.Find(bson.M{"email": email})

	var user model.User
	err := query.One(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// ===== Cart =====

// Inserts the
func (db *Database) PutCart(cart *model.Cart) error {

	id := cart.Id
	if id == "" {
		id = bson.NewObjectId()
	}

	info, err := db.carts.Upsert(bson.M{"_id": id}, cart)
	if err != nil {
		return err
	}

	newId, assertionOk := info.UpsertedId.(bson.ObjectId)
	if assertionOk {
		cart.Id = newId
	}

	return nil
}

func (db *Database) GetCart(cartId string) (model.Cart, error) {

	query := db.carts.Find(bson.M{"_id": bson.ObjectIdHex(cartId)})

	var cart model.Cart
	err := query.One(&cart)
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (db *Database) ValidateCart(userId bson.ObjectId, cartId string, quantities map[string]int) (*model.Cart, error) {

	ids := make([]string, len(quantities))
	for id, _ := range quantities {
		ids = append(ids, id)
	}

	query := db.products.Find(bson.M{"id": bson.M{"$in": ids}})

	var products []model.Product
	err := query.All(&products)
	if err != nil {
		return nil, err
	}

	total := 0
	cartItems := make([]model.CartItem, 0)

	for _, p := range products {
		quantity := quantities[p.Id]
		total += p.Price * quantity
		cartItems = append(cartItems, model.CartItem{Product: p, Quantity: quantity})
	}

	var id bson.ObjectId
	if cartId == "" {
		id = bson.ObjectId("")
	} else {
		id = bson.ObjectIdHex(cartId)
	}

	cart := model.Cart{
		Items:    cartItems,
		Expires:  time.Now().Add(1 * time.Hour).Unix(),
		Total:    total,
		Shipping: 599,
		UserId:   userId,
		Id:       id,
	}

	err = db.PutCart(&cart)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	return &cart, nil
}

// ===== Order =====

func (db *Database) PutOrder(order *model.Order) error {

	info, err := db.orders.Upsert(bson.M{"_id": bson.NewObjectId()}, order)
	if err != nil {
		return err
	}

	newId, assertionOk := info.UpsertedId.(bson.ObjectId)
	if assertionOk {
		order.Id = newId
	}

	return nil
}

// ===== User =====

func (db *Database) NewUser(email string, password string) (model.User, error) {

	query := db.users.Find(bson.M{"email": email})

	var existingUser model.User
	err := query.One(&existingUser)
	if err != nil {
		if err != mgo.ErrNotFound {
			return model.User{}, err
		}
	} else {
		// no error => user already exists
		return model.User{}, fmt.Errorf("user with email: %s already exists", email)
	}

	hashedPassword, err := argon.Hash(password)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{Email: email, Password: hashedPassword}
	err = db.users.Insert(user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// ===== Helper =====

func getProductsResults(query *mgo.Query) ([]model.Product, error) {

	products := make([]model.Product, 0)
	err := query.All(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (db *Database) GetNthPage(query bson.M, page int) ([]model.Product, bool, error) {

	if page < 1 {
		return nil, false, fmt.Errorf("page shoud be between 1-n, got: %d", page)
	}

	q := db.products.Find(query).Skip((page - 1) * pageSize).Limit(pageSize)

	results, err := getProductsResults(q)
	if err != nil {
		return results, false, err
	}

	if len(results) < pageSize {
		isLast, err := db.isLastPage(query, page+1)
		if err != nil {
			return results, false, err
		}

		if isLast {
			return results, true, nil
		}
	}

	return results, false, nil
}

func (db *Database) isLastPage(query bson.M, potentialLastPage int) (bool, error) {
	q := db.products.Find(query).Skip((potentialLastPage - 1) * pageSize).Limit(pageSize)

	results, err := getProductsResults(q)
	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return true, nil
	}

	return false, nil
}
