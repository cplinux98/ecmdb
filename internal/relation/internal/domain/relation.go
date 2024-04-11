package domain

import (
	"time"
)

const (
	MappingOneToOne   = iota + 1 // 一对一关系
	MappingOneToMany             // 一对多关系
	MappingManyToMany            // 多对多关系
)

type ModelRelation struct {
	ID                     int64
	SourceModelIdentifies  string
	TargetModelIdentifies  string
	RelationTypeIdentifies string // 关联类型唯一索引
	RelationName           string // 拼接字符
	Mapping                string // 关联关系
	Ctime                  time.Time
	Utime                  time.Time
}

type ResourceRelation struct {
	ID                     int64
	SourceModelIdentifies  string
	TargetModelIdentifies  string
	SourceResourceID       string
	TargetResourceID       string
	RelationTypeIdentifies string // 关联类型唯一索引
	RelationName           string // 拼接字符
}