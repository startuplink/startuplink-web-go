package app

import "github.com/pkg/errors"

type Profile int

const (
	LOCAL = 0
	PROD  = 1
)

var names = [...]string{
	"local",
	"prod",
}

func (profile Profile) String() string {
	if profile < LOCAL || profile > PROD {
		return "Unknown"
	}
	return names[profile]
}

func getProfile(name string) (Profile, error) {
	switch name {
	case "local":
		return LOCAL, nil
	case "prod":
		return PROD, nil
	default:
		return -1, errors.New("profile not exists")
	}
}
