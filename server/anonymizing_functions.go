package server

import (
	"crypto/sha1"
	"encoding/base64"
)

func hashString(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

type AnonymizingFunction func(interface{}) interface{}

func email(original interface{}) interface{} {
	return hashString(original.(string))
}

func passwordStars(original interface{}) interface{} {
	return "********"
}

func nameFirstDwarf(original interface{}) interface{} {
	return "Sleepy"
}
func nameMiddleDwarf(original interface{}) interface{} {
	return "Grumpy"
}
func nameLastDwarf(original interface{}) interface{} {
	return "Dopey"
}
