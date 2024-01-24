package vrule

import "github.com/protomem/gotube/pkg/validation"

func Nickname(v *validation.Validator, nickname string) {
	v.CheckField(validation.MinRunes(nickname, 4), "nickname", "must be at least 4 characters long")
	v.CheckField(validation.MaxRunes(nickname, 32), "nickname", "must be at most 32 characters long")
}

func Password(v *validation.Validator, password string) {
	v.CheckField(validation.MinRunes(password, 6), "password", "must be at least 6 characters long")
	v.CheckField(validation.MaxRunes(password, 18), "password", "must be at most 18 characters long")
}

func Email(v *validation.Validator, email string) {
	v.CheckField(validation.IsEmail(email), "email", "must be a valid email address")
}

func AvatarPath(v *validation.Validator, avatarPath string) {
	v.CheckField(validation.IsPath(avatarPath), "avatarPath", "must be a valid path")
}

func Description(v *validation.Validator, description string) {
	v.CheckField(validation.MaxRunes(description, 500), "description", "must be at most 500 characters long")
}
