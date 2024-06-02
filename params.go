package pbclient

import (
	"fmt"
	"net/url"
	"strings"
)

type Params struct {
	// Auth token. Required for all operations except for accessing public files.
	Token string

	// The name of the collection to be searched
	Collection string

	// Collection data (ONLY applies to Create/Update operations)
	Data interface{}

	// Item ID (ONLY applies to View/Delete operations)
	ID string

	// FileName (ONLY applies when fetching files)
	FileName string

	// The page (aka. offset) of the paginated list (default to 1).
	Page int

	// Specify the max returned records per page (default to 30).
	PerPage int

	// Specify the records order attribute(s).
	// Add - / + (default) in front of the attribute for DESC / ASC order. Ex.:
	// DESC by created and ASC by id: ?sort=-created,id
	Sort string

	// Filter the returned records. Ex.:
	// ?filter=(id='abc' && created>'2022-01-01')
	// The syntax basically follows the format OPERAND OPERATOR OPERAND, where:
	// OPERAND - could be any of the above field literal, string (single or double quoted), number, null, true, false
	// OPERATOR - is one of:
	// 		= Equal
	// 		!= NOT equal
	// 		> Greater than
	// 		>= Greater than or equal
	// 		< Less than
	// 		<= Less than or equal
	// 		~ Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
	// 		!~ NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
	// 		?= Any/At least one of Equal
	// 		?!= Any/At least one of NOT equal
	// 		?> Any/At least one of Greater than
	// 		?>= Any/At least one of Greater than or equal
	// 		?< Any/At least one of Less than
	// 		?<= Any/At least one of Less than or equal
	// 		?~ Any/At least one of Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
	// 		?!~ Any/At least one of NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
	// To group and combine several expressions you could use brackets (...), && (AND) and || (OR) tokens.
	Filter string

	// Auto expand record relations. Ex.:
	// ?expand=relField1,relField2.subRelField
	// Supports up to 6-levels depth nested relations expansion.
	// The expanded relations will be appended to each individual record under the expand property (eg. "expand": {"relField1": {...}, ...}).
	// Only the relations to which the request user has permissions to view will be expanded.
	Expand string

	// Comma separated string of the fields to return in the JSON response (by default returns all fields). Ex.:
	// ?fields=*,expand.relField.name
	// * targets all keys from the specific depth level.
	// In addition, the following field modifiers are also supported:
	// :excerpt(maxLength, withEllipsis?)
	// Returns a short plain text version of the field string value.
	// Ex.: ?fields=*,description:excerpt(200,true)
	Fields string

	// If it is set the total counts query will be skipped and the response fields totalItems and totalPages will have -1 value.
	// This could drastically speed up the search queries when the total counters are not needed or cursor based pagination is used.
	// For optimization purposes, it is set by default for the getFirstListItem() and getFullList() SDKs methods.
	SkipTotal bool

	// When requesting image files thumb configuration can be specified. The following thumb formats are currently supported:
	// - WxH (eg. 100x300) - crop to WxH viewbox (from center)
	// - WxHt (eg. 100x300t) - crop to WxH viewbox (from top)
	// - WxHb (eg. 100x300b) - crop to WxH viewbox (from bottom)
	// - WxHf (eg. 100x300f) - fit inside a WxH viewbox (without cropping)
	// - 0xH (eg. 0x300) - resize to H height preserving the aspect ratio
	// - Wx0 (eg. 100x0) - resize to W width preserving the aspect ratio
	Thumb string
}

func (p Params) QueryString() string {
	var sb strings.Builder

	if p.Page > 0 {
		sb.WriteString(fmt.Sprintf("&page=%d", p.Page))
	}
	if p.PerPage > 0 {
		sb.WriteString(fmt.Sprintf("&perPage=%d", p.PerPage))
	}
	if p.Sort != "" {
		sb.WriteString(fmt.Sprintf("&sort=%s", p.Sort))
	}
	if p.Filter != "" {
		sb.WriteString(fmt.Sprintf("&filter=%s", url.QueryEscape(p.Filter)))
	}
	if p.Expand != "" {
		sb.WriteString(fmt.Sprintf("&expand=%s", p.Expand))
	}
	if p.Fields != "" {
		sb.WriteString(fmt.Sprintf("&fields=%s", p.Expand))
	}
	if p.SkipTotal {
		sb.WriteString("&skipTotal=true")
	}

	return strings.TrimPrefix(sb.String(), "&")
}
