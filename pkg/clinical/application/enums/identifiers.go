package enums

import (
	"fmt"
	"io"
	"strconv"
)

// IDDocumentType is an internal code system for identification document types.
type IDDocumentType string

// ID type constants
const (
	// IDDocumentTypeNationalID is the national id number of the patient
	IDDocumentTypeNationalID IDDocumentType = "national_id"
	// IDDocumentTypePassport represents patient identification of type passport
	IDDocumentTypePassport IDDocumentType = "passport"
	// IDDocumentTypeAlienID represents the alien id used to identify an patient of alien origin
	IDDocumentTypeAlienID IDDocumentType = "alien_id"
	// IDDocumentTypeCCC is represents the CCC number used to identify a patient
	IDDocumentTypeCCC IDDocumentType = "ccc_number"
)

// IsValid checks that the ID type is valid
func (e IDDocumentType) IsValid() bool {
	switch e {
	case IDDocumentTypeNationalID, IDDocumentTypePassport, IDDocumentTypeAlienID, IDDocumentTypeCCC:
		return true
	}

	return false
}

// String ...
func (e IDDocumentType) String() string {
	return string(e)
}

// UnmarshalGQL translates the input value to an ID type
func (e *IDDocumentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IDDocumentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IDDocumentType", str)
	}

	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IDDocumentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
