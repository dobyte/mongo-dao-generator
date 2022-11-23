package template

const ExternalTemplate = `
package ${VarDaoPackageName}

import (
	"${VarDaoPackagePath}/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type ${VarDaoPrefixName}Columns = internal.${VarDaoPrefixName}Columns

type ${VarDaoClassName} struct {
	*internal.${VarDaoClassName}
}

func New${VarDaoClassName}(db *mongo.Database) *${VarDaoClassName} {
	return &${VarDaoClassName}{${VarDaoClassName}: internal.New${VarDaoClassName}(db)}
}
`

const InternalTemplate = `
// --------------------------------------------------------------------------------------------
// The following code is automatically generated by the gen-mongo-dao tool. 
// Please do not modify this code manually to avoid being overwritten in the next generation. 
// For more tool details, please click the link to view https://github.com/dobyte/gen-mongo-dao
// --------------------------------------------------------------------------------------------

package internal

import (
	${VarPackages}
)

type ${VarDaoPrefixName}FilterFunc func(cols *${VarDaoPrefixName}Columns) interface{}
type ${VarDaoPrefixName}UpdateFunc func(cols *${VarDaoPrefixName}Columns) interface{}
type ${VarDaoPrefixName}FindOneOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.FindOneOptions
type ${VarDaoPrefixName}FindManyOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.FindOptions
type ${VarDaoPrefixName}UpdateOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.UpdateOptions
type ${VarDaoPrefixName}DeleteOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.DeleteOptions
type ${VarDaoPrefixName}InsertOneOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.InsertOneOptions
type ${VarDaoPrefixName}InsertManyOptionsFunc func(cols *${VarDaoPrefixName}Columns) *options.InsertManyOptions

type ${VarDaoClassName} struct {
	Columns    *${VarDaoPrefixName}Columns
	Database   *mongo.Database
	Collection *mongo.Collection
}

type ${VarDaoPrefixName}Columns struct {
	${VarModelColumnsDefine}
}

var ${VarDaoVariableName}Columns = &${VarDaoPrefixName}Columns{
	${VarModelColumnsInstance}
}

func New${VarDaoClassName}(db *mongo.Database) *${VarDaoClassName} {
	return &${VarDaoClassName}{
		Columns:    ${VarDaoVariableName}Columns,
		Database:   db,
		Collection: db.Collection("${VarCollectionName}"),
	}
}

// InsertOne executes an insert command to insert a single document into the collection.
func (dao *${VarDaoClassName}) InsertOne(ctx context.Context, model *${VarModelPackageName}.${VarModelClassName}, optionsFunc ...${VarDaoPrefixName}InsertOneOptionsFunc) (*mongo.InsertOneResult, error) {
	if model == nil {
		return nil, errors.New("model is nil")
	}

	if err := dao.autofill(ctx, model); err != nil {
		return nil, err
	}

	var opts *options.InsertOneOptions

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.InsertOne(ctx, model, opts)
}

// InsertMany executes an insert command to insert multiple documents into the collection.
func (dao *${VarDaoClassName}) InsertMany(ctx context.Context, models []*${VarModelPackageName}.${VarModelClassName}, optionsFunc ...${VarDaoPrefixName}InsertManyOptionsFunc) (*mongo.InsertManyResult, error) {
	if len(models) == 0 {
		return nil, errors.New("models is empty")
	}

	documents := make([]interface{}, 0, len(models))
	for i := range models {
		model := models[i]
		if err := dao.autofill(ctx, model); err != nil {
			return nil, err
		}
		documents = append(documents, model)
	}

	var opts *options.InsertManyOptions

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.InsertMany(ctx, documents, opts)
}

// UpdateOne executes an update command to update at most one document in the collection.
func (dao *${VarDaoClassName}) UpdateOne(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, updateFunc ${VarDaoPrefixName}UpdateFunc, optionsFunc ...${VarDaoPrefixName}UpdateOptionsFunc) (*mongo.UpdateResult, error) {
	var (
		opts   *options.UpdateOptions
		filter = filterFunc(dao.Columns)
		update = updateFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.UpdateOne(ctx, filter, update, opts)
}

// UpdateOneByID executes an update command to update at most one document in the collection.
func (dao *${VarDaoClassName}) UpdateOneByID(ctx context.Context, id string, updateFunc ${VarDaoPrefixName}UpdateFunc, optionsFunc ...${VarDaoPrefixName}UpdateOptionsFunc) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return dao.UpdateOne(ctx, func(cols *Columns) interface{} {
		return bson.M{"_id": objectID}
	}, updateFunc, optionsFunc...)
}

// UpdateMany executes an update command to update documents in the collection.
func (dao *${VarDaoClassName}) UpdateMany(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, updateFunc ${VarDaoPrefixName}UpdateFunc, optionsFunc ...${VarDaoPrefixName}UpdateOptionsFunc) (*mongo.UpdateResult, error) {
	var (
		opts   *options.UpdateOptions
		filter = filterFunc(dao.Columns)
		update = updateFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.UpdateMany(ctx, filter, update, opts)
}

// FindOne executes a find command and returns a model for one document in the collection.
func (dao *${VarDaoClassName}) FindOne(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, optionsFunc ...${VarDaoPrefixName}FindOneOptionsFunc) (*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		opts   *options.FindOneOptions
		model  = &${VarModelPackageName}.${VarModelClassName}{}
		filter = filterFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	err := dao.Collection.FindOne(ctx, filter, opts).Decode(model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

// FindOneByID executes a find command and returns a model for one document in the collection.
func (dao *${VarDaoClassName}) FindOneByID(ctx context.Context, id string, optionsFunc ...${VarDaoPrefixName}FindOneOptionsFunc) (*${VarModelPackageName}.${VarModelClassName}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return dao.FindOne(ctx, func(cols *Columns) interface{} {
		return bson.M{"_id": objectID}
	}, optionsFunc...)
}

// FindMany executes a find command and returns many models the matching documents in the collection.
func (dao *${VarDaoClassName}) FindMany(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, optionsFunc ...${VarDaoPrefixName}FindManyOptionsFunc) ([]*${VarModelPackageName}.${VarModelClassName}, error) {
	var (
		opts   *options.FindOptions
		filter = filterFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	cur, err := dao.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	models := make([]*${VarModelPackageName}.${VarModelClassName}, 0)
	
	if err = cur.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// DeleteOne executes a delete command to delete at most one document from the collection.
func (dao *${VarDaoClassName}) DeleteOne(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, optionsFunc ...${VarDaoPrefixName}DeleteOptionsFunc) (*mongo.DeleteResult, error) {
	var (
		opts   *options.DeleteOptions
		filter = filterFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.DeleteOne(ctx, filter, opts)
}

// DeleteOneByID executes a delete command to delete at most one document from the collection.
func (dao *${VarDaoClassName}) DeleteOneByID(ctx context.Context, id string, optionsFunc ...${VarDaoPrefixName}DeleteOptionsFunc) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return dao.DeleteOne(ctx, func(cols *Columns) interface{} {
		return bson.M{"_id": objectID}
	}, optionsFunc...)
}

// DeleteMany executes a delete command to delete documents from the collection.
func (dao *${VarDaoClassName}) DeleteMany(ctx context.Context, filterFunc ${VarDaoPrefixName}FilterFunc, optionsFunc ...${VarDaoPrefixName}DeleteOptionsFunc) (*mongo.DeleteResult, error) {
	var (
		opts   *options.DeleteOptions
		filter = filterFunc(dao.Columns)
	)

	if len(optionsFunc) > 0 {
		opts = optionsFunc[0](dao.Columns)
	}

	return dao.Collection.DeleteMany(ctx, filter, opts)
}

// autofill when inserting data
func (dao *${VarDaoClassName}) autofill(ctx context.Context, model *${VarModelPackageName}.${VarModelClassName}) error {
	${VarAutofillCode}
}
`
