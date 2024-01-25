package handler

type Video struct {
	*Base
}

func NewVideo() *Video {
	return &Video{NewBase()}
}
