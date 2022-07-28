package model

import "time"

type Client struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`

	UUID string `json:"uuid" xorm:"'uuid' index"`
	SN   string `json:"sn" xorm:"'sn' index"`
	CPU  string `json:"cpu" xorm:"'cpu' index"`
	MAC  string `json:"mac" xorm:"'mac' index"`

	Disabled bool      `json:"disabled"`
	Updated  time.Time `json:"updated" xorm:"updated"`
	Created  time.Time `json:"created" xorm:"created"`
}
