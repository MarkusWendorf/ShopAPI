package database

import "gopkg.in/mgo.v2/bson"

type QueryBuilder struct {
	constraints []bson.M
}

func (builder *QueryBuilder) Category(category string) *QueryBuilder {

	if category == "" {
		return builder
	}

	constraint := bson.M{"category": category}

	builder.constraints = append(builder.constraints, constraint)
	return builder
}

func (builder *QueryBuilder) Name(name string) *QueryBuilder {

	if name == "" {
		return builder
	}

	constraint := bson.M{"pname": bson.M{"$regex": name, "$options": "i"}}

	builder.constraints = append(builder.constraints, constraint)
	return builder
}

func (builder *QueryBuilder) Price(from int, to int) *QueryBuilder {

	constraint := bson.M{"price": bson.M{"$gte": from, "$lte": to}}

	builder.constraints = append(builder.constraints, constraint)
	return builder
}

func (builder *QueryBuilder) Build() bson.M {
	return bson.M{"$and": builder.constraints}
}
