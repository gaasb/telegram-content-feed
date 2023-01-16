package bot

type UserType uint

const (
	PostEdit UserType = iota
	PostDelete
	CreateTag
)

const (
	RoleAdmin     = "admin"
	RoleModerator = "mod"
)

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
