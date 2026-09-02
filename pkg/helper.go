package pkg

import (
	"time"

	"github.com/goombaio/namegenerator"
)

func GenerateUniqueName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()
	return name
}
