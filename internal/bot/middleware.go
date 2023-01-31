package bot

import (
	"errors"
	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

type Administrator struct {
	UserID int64      `bson:"_id" json:"_id"`
	Rights []UserType `bson:"rights" json:"rights"`
}
type AdministratorRoleAndRights struct {
	Rights []UserType `bson:"rights" json:"rights"`
}
type UserType uint
type UserRole string

const (
	PostEdit UserType = iota
	PostDelete
	CreateTag
)

var administrators = map[int64][]UserType{}

const (
	RoleAdmin     = "admin"
	RoleModerator = "mod"
)

var adminOnlyController *telebot.Group

var userRights = map[UserType]string{
	PostEdit:   "Edit message",
	PostDelete: "Delete message",
	CreateTag:  "Create tags",
}
var rosterList map[string]*Role

type Role struct {
	Action
}

func (r *Role) RoleByName(role string) Action {
	switch role {
	case RoleAdmin:
		return &Admin{}
	case RoleModerator:
		return &Moderator{}
	default:
		return nil
	}
}

type Action interface {
	Do()
}

type Admin Role
type Moderator Role

func (a *Admin) Do() {
}
func (a *Moderator) Do() {
}

func UpdateRosterList() {
}
func AddRoleToUser(user interface{}, role string) {
	_ = Role{}
	//rosterList = append(rosterList, user, r.RoleByName(role))
}

func ForAdministrators(right UserType) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			if admin, ok := administrators[ctx.Sender().ID]; ok && isContainRight(admin, right) {
				return next(ctx)
			}
			return nil
		}

	}
}
func AddAdmin(userID int64, rights ...UserType) error {
	if err := InsertAdministrator(&Administrator{UserID: userID, Rights: rights}); err != nil {
		return errors.New("Already in db")
	}
	administrators[userID] = rights
	return nil
}
func InitAllAdministrators() {
	administrators = FindAllAdministrators()

}
func isContainRight(role []UserType, right UserType) bool {
	return slices.Contains(role, right)
}
