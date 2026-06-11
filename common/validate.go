package common

import (
	"github.com/go-playground/validator/v10"
)

// Validate is the global validator instance
var Validate = validator.New()

