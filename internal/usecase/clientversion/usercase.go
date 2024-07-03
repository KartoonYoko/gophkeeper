package clientversion

import "fmt"

type Usecase struct {
	version   string
	buildDate string
}

func New(buildVersion string, buildDate string) *Usecase {
	return &Usecase{
		version:   buildVersion,
		buildDate: buildDate,
	}
}

func (u *Usecase) Version() string {
	return u.version
}

func (u *Usecase) BuildDate() string {
	return u.buildDate
}

func (u *Usecase) Info() string {
	return fmt.Sprintf("Version: %s\nBuild date: %s", u.Version(), u.BuildDate())
}
