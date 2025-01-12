package dao

import (
	"context"
	"fmt"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UserCollection = "c_user"

type UserDAO interface {
	CreatUser(ctx context.Context, user User) (int64, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	ListUser(ctx context.Context, offset, limit int64) ([]User, error)
	AddOrUpdateRoleBind(ctx context.Context, id int64, roleCodes []string) (int64, error)
	Count(ctx context.Context) (int64, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
}

type userDao struct {
	db *mongox.Mongo
}

func (dao *userDao) UpdatePassword(ctx context.Context, id int64, password string) error {
	col := dao.db.Collection(UserCollection)
	updateDoc := bson.M{
		"$set": bson.M{
			"password": password,
			"utime":    time.Now().UnixMilli(),
		},
	}

	filter := bson.M{"id": id}
	_, err := col.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("修改文档操作: %w", err)
	}

	return nil
}

func (dao *userDao) FindById(ctx context.Context, id int64) (User, error) {
	col := dao.db.Collection(UserCollection)
	var u User
	filter := bson.M{"id": id}

	if err := col.FindOne(ctx, filter).Decode(&u); err != nil {
		return User{}, fmt.Errorf("解码错误，%w", err)
	}

	return u, nil
}

func (dao *userDao) AddOrUpdateRoleBind(ctx context.Context, id int64, roleCodes []string) (int64, error) {
	col := dao.db.Collection(UserCollection)
	updateDoc := bson.M{
		"$set": bson.M{
			"role_codes": roleCodes,
			"utime":      time.Now().UnixMilli(),
		},
	}
	filter := bson.M{"id": id}
	count, err := col.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return 0, fmt.Errorf("修改文档操作: %w", err)
	}

	return count.ModifiedCount, nil
}

func (dao *userDao) ListUser(ctx context.Context, offset, limit int64) ([]User, error) {
	col := dao.db.Collection(UserCollection)
	filter := bson.M{}
	opts := &options.FindOptions{
		Sort:  bson.D{{Key: "ctime", Value: -1}},
		Limit: &limit,
		Skip:  &offset,
	}

	cursor, err := col.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询错误, %w", err)
	}

	var result []User
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解码错误: %w", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("游标遍历错误: %w", err)
	}
	return result, nil
}

func (dao *userDao) Count(ctx context.Context) (int64, error) {
	col := dao.db.Collection(UserCollection)
	filer := bson.M{}

	count, err := col.CountDocuments(ctx, filer)
	if err != nil {
		return 0, fmt.Errorf("文档计数错误: %w", err)
	}

	return count, nil
}

func NewUserDao(db *mongox.Mongo) UserDAO {
	return &userDao{
		db: db,
	}
}

func (dao *userDao) CreatUser(ctx context.Context, user User) (int64, error) {
	now := time.Now()
	user.Ctime, user.Utime = now.UnixMilli(), now.UnixMilli()
	user.Id = dao.db.GetIdGenerator(UserCollection)
	col := dao.db.Collection(UserCollection)

	_, err := col.InsertOne(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("插入数据错误: %w", err)
	}

	return user.Id, nil
}

func (dao *userDao) FindByUsername(ctx context.Context, username string) (User, error) {
	col := dao.db.Collection(UserCollection)
	var u User
	filter := bson.M{"username": username}

	if err := col.FindOne(ctx, filter).Decode(&u); err != nil {
		return User{}, fmt.Errorf("解码错误，%w", err)
	}

	return u, nil
}

type User struct {
	Id          int64    `bson:"id"`
	Username    string   `bson:"username"`
	Password    string   `bson:"password"`
	Email       string   `bson:"email"`
	Title       string   `bson:"title"`
	DisplayName string   `bson:"display_name"`
	CreateType  uint8    `bson:"create_type"`
	Status      uint8    `bson:"status"`
	Ctime       int64    `bson:"ctime"`
	Utime       int64    `bson:"utime"`
	RoleCodes   []string `bson:"role_codes"`
}
