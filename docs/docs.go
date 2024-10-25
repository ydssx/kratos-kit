package docs

import (
	_ "embed"
	"encoding/json"
)

//go:embed apidocs.swagger.json
var ApiDocs json.RawMessage
