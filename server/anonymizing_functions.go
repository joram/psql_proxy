package server

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/goombaio/namegenerator"
	"math/rand"
	"strings"
)

func hashString(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func getSeed(s string) int64 {
	var myint int64
	hashBytes := sha1.Sum([]byte(s))
	buf := bytes.NewBuffer(hashBytes[:])
	binary.Read(buf, binary.LittleEndian, &myint)
	return myint
}

var anonymizingFunctions = map[string]AnonymizingFunction{
	"email":          email,
	"password.stars": passwordStars,
	"name.first":     name,
	"name.middle":    name,
	"name.last":      name,
}

type AnonymizingFunction func(interface{}) interface{}

func email(original interface{}) interface{} {
	parts := strings.SplitN(fmt.Sprintf("%s", original), "@", 2)
	username := parts[0]
	domain := parts[1]
	hashedEmail := hashString(username) + "@" + hashString(domain) + ".com"
	return hashedEmail
}

func passwordStars(_ interface{}) interface{} {
	// returns a string of 5 to 15 stars
	n := 5 + rand.Intn(10)
	return strings.Repeat("*", n)
}

func name(original interface{}) interface{} {
	seed := getSeed(original.(string))
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()
	return name
}
