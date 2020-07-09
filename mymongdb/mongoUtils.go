package mymongdb

import (
	"context"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//插入单个文档
func (m *mgo) InsertOne(value interface{}) *mongo.InsertOneResult {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		mylog.Error("InsertOne error:%s", err)
	}
	return insertResult
}

//插入多个文档
func (m *mgo) InsertMany(values []interface{}) *mongo.InsertManyResult {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.InsertMany(context.TODO(), values)
	if err != nil {
		mylog.Error("InsertMany error:%s", err)
	}
	return result
}

//删除文档
func (m *mgo) Delete(key string, value interface{}) int64 {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		mylog.Error("Delete error:%s", err)
	}
	return count.DeletedCount
}

//删除多个文档
func (m *mgo) DeleteMany(key string, value interface{}) int64 {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	filter := bson.D{{key, value}}

	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		mylog.Error("DeleteMany error:%s", err)
	}
	return count.DeletedCount
}

//更新单个文档
func (m *mgo) UpdateOne(filter, update interface{}) int64 {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("UpdateOne error:%s", err)
	}
	return result.UpsertedCount
}

//更新多个文档
func (m *mgo) UpdateMany(filter, update interface{}) int64 {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("UpdateMany error:%s", err)
	}
	return result.UpsertedCount
}

// 查询单个文档
func (m *mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	client := DB.Mongo
	collection, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		mylog.Error("FindOne error:%s", err)
	}
	filter := bson.D{{key, value}}
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

//查询多个文档
func (m *mgo) FindMany(filter interface{}) (*mongo.Cursor, error) {
	client := DB.Mongo
	collection, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		mylog.Error("FindMany error:%s", err)
	}
	return collection.Find(context.TODO(), filter)
}

//多条件查询
func (m *mgo) FindManyByFilters(filter interface{}) (*mongo.Cursor, error) {
	client := DB.Mongo
	collection, err := client.Database(m.database).Collection(m.collection).Clone()
	if err != nil {
		mylog.Error("FindMany error:%s", err)
	}
	return collection.Find(context.TODO(), bson.M{"$and": filter})
}

//查询集合里有多少数据
func (m *mgo) CollectionCount() (string, int64) {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

//按选项查询集合
// Skip 跳过
// Limit 读取数量
// sort 1 ，-1 . 1 为升序 ， -1 为降序
func (m *mgo) CollectionDocuments(Skip, Limit int64, sort int, key string, value interface{}) *mongo.Cursor {
	client := DB.Mongo
	collection := client.Database(m.database).Collection(m.collection)
	SORT := bson.D{{"_id", sort}}
	filter := bson.D{{key, value}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	temp, _ := collection.Find(context.Background(), filter, findOptions)
	return temp
}
