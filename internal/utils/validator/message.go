package validator

var Messages = map[string]string{
	// ===== Common =====
	"required": "is required",

	// ===== String =====
	"min": "must be at least {param} characters",
	"max": "must not exceed {param} characters",

	// ===== Number =====
	"gt":  "must be greater than {param}",
	"gte": "must be greater than or equal to {param}",
	"lt":  "must be less than {param}",
	"lte": "must be less than or equal to {param}",
	"eq":  "must be equal to {param}",
	"ne":  "must not be equal to {param}",

	// ===== Format =====
	"email": "must be a valid email address",
	"url":   "must be a valid URL",
	"uuid":  "must be a valid UUID",

	// ===== Boolean =====
	"boolean": "must be a boolean value",

	// ===== Slice / Array =====
	"unique": "elements must be unique",

	// ===== Custom =====
}
