package utility

import (
	"log"

	"github.com/pkg/errors"
)

func MustCheckError(err error) {
	if err != nil {
		wrapErr := errors.Wrap(err, "")
		log.Fatalf("%+v\n", wrapErr)
	}
}
