//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Users struct {
	ID        string `sql:"primary_key"`
	Name      string
	Username  string
	Password  []byte
	CreatorID string
	CreatedAt time.Time
	Role      UsersRole
}