package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
)

const (
	patientLinkExpireMinutes = 30
	maxAccessAttempts        = 5
)

// PatientLink stores a map of patient IDs to short lived opaque IDs.
//
// These opaque IDs are used in publicly visible links.
// The intention is to obscure confidential (long lived) patient IDs.
type PatientLink struct {
	ID        string    `json:"ID" firestore:"ID"`
	PatientID string    `json:"patientID" firestore:"patientID"`
	OpaqueID  string    `json:"opaqueID" firestore:"opaqueID"`
	Expires   time.Time `json:"expires" firestore:"expires"`
	Deleted   bool      `json:"deleted" firestore:"deleted"`
}

// PatientLinkConnection is used to serialize GraphQL Relay connections for patient links
type PatientLinkConnection struct {
	Edges    []*PatientLinkEdge      `json:"edges"`
	PageInfo *firebasetools.PageInfo `json:"pageInfo"`
}

// PatientLinkEdge is used to serialize GraphQL relay edges for patient links
type PatientLinkEdge struct {
	Cursor *string      `json:"cursor"`
	Node   *PatientLink `json:"node"`
}

// IsNode marks this struct as a relay Node
func (pl *PatientLink) IsNode() {}

// GetID returns the patient links primary key
func (pl *PatientLink) GetID() firebasetools.ID {
	return firebasetools.IDValue(pl.ID)
}

// SetID sets the patient links' id
func (pl *PatientLink) SetID(id string) {
	pl.ID = id
}

// GeneratePatientLink creates a random patient link
func GeneratePatientLink(patientID string) (*PatientLink, error) {
	xid := xid.New()
	ts := time.Now().Add(time.Minute * patientLinkExpireMinutes)
	ctx := context.Background()
	link := &PatientLink{
		PatientID: patientID,
		OpaqueID:  xid.String(),
		Expires:   ts,
	}
	_, _, err := firebasetools.CreateNode(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("unable to save patient link: %w", err)
	}
	return link, nil
}

// GetPatientID returns the actual patient ID
func GetPatientID(ctx context.Context, opaqueID string) (string, error) {
	filterParam := firebasetools.FilterParam{
		FieldName:           "opaqueID",
		FieldType:           enumutils.FieldTypeString,
		ComparisonOperation: enumutils.OperationEqual,
		FieldValue:          opaqueID,
	}

	filterParamExpires := firebasetools.FilterParam{
		FieldName:           "expires",
		FieldType:           enumutils.FieldTypeTimestamp,
		ComparisonOperation: enumutils.OperationGreaterThanOrEqualTo,
		FieldValue:          time.Now(),
	}

	filter := firebasetools.FilterInput{
		FilterBy: []*firebasetools.FilterParam{
			&filterParam,
			&filterParamExpires,
		},
	}
	docs, _, err := firebasetools.QueryNodes(ctx, nil, &filter, nil, &PatientLink{})
	if err != nil {
		return "", err
	}

	if len(docs) < 1 {
		return "", fmt.Errorf("expected 1 matching id, got %d", len(docs))
	}
	doc := docs[0]
	var pl PatientLink
	err = doc.DataTo(&pl)
	if err != nil {
		return "", fmt.Errorf("unable to read PatientLink from firebase snapshot: %w", err)
	}
	docmap := doc.Data()
	patientID := docmap["patientID"]
	patientIDstr, ok := patientID.(string)
	if !ok {
		return "", fmt.Errorf("id is not a string, it is of type %T", patientID)
	}

	return patientIDstr, err
}
