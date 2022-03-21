package infrastructure

import (
	"context"
	"log"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/firebasetools"
)

// FHIRRepository ...
type FHIRRepository interface {
	CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error)
	DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	GetFHIRPatientAllData(fhirResourceID string) ([]byte, error)
	FHIRRestURL() string
	GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	GetFHIRPatientEverything(fhirResourceID string) ([]byte, error)
}

// FHIRService is an implementation of the database repository
// It is implementation agnostic i.e logic should be handled using
// the preferred database
type FHIRService struct {
	FHIR FHIRRepository
}

// NewFHIRService creates a new database service
func NewFHIRService() FHIRService {
	repo := fhir.NewFHIRRepository()

	return FHIRService{
		repo,
	}
}

// CreateFHIRResource creates an FHIR resource.
//
// The payload should be the result of marshalling a resource to JSON
func (d FHIRService) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return d.FHIR.CreateFHIRResource(resourceType, payload)
}

// DeleteFHIRResource deletes an FHIR resource.
func (d FHIRService) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return d.FHIR.DeleteFHIRResource(resourceType, fhirResourceID)
}

// PatchFHIRResource patches a FHIR resource.
// The payload is a JSON patch document that follows guidance on Patch from the
// FHIR standard.
// See:
// payload := []map[string]interface{}{
// 	{
// 		"op":    "replace",
// 		"path":  "/active",
// 		"value": active,
// 	},
// }
// See: https://www.hl7.org/fhir/http.html#patch
func (d FHIRService) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return d.FHIR.PatchFHIRResource(resourceType, fhirResourceID, payload)
}

// UpdateFHIRResource updates the entire contents of a resource.
func (d FHIRService) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return d.FHIR.UpdateFHIRResource(resourceType, fhirResourceID, payload)
}

// GetFHIRPatientAllData gets all resources associated with a particular
// patient compartment.
func (d FHIRService) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	return d.FHIR.GetFHIRPatientAllData(fhir.DatasetLocation)
}

// FHIRRestURL composes a FHIR REST URL for manual calls to the FHIR REST API
func (d FHIRService) FHIRRestURL() string {
	return d.FHIR.FHIRRestURL()
}

// GetFHIRResource gets an FHIR resource.
func (d FHIRService) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return d.FHIR.GetFHIRResource(resourceType, fhirResourceID)
}

// GetFHIRPatientEverything gets all fhir patient resource information.
func (d FHIRService) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	return d.FHIR.GetFHIRPatientEverything(fhirResourceID)
}

// Repository ...
type Repository interface {
	SaveEmailOTP(ctx context.Context, email string, optIn bool) error
	StageStartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error
}

// DBService is an implementation of the database repository
// It is implementation agnostic i.e logic should be handled using
// the preferred database
type DBService struct {
	firestore *fb.Repository
}

// NewDBService creates a new database service
func NewDBService() *DBService {
	ctx := context.Background()
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil
	}
	fbc, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}
	fsc, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}
	firestoreExtension := fb.NewFirestoreClientExtension(fsc)

	firestore := fb.NewFirebaseRepository(firestoreExtension, fbc)
	return &DBService{
		firestore: firestore,
	}
}

// SaveEmailOTP  persist the data of the newly created OTP to a datastore
func (db DBService) SaveEmailOTP(
	ctx context.Context,
	email string, optIn bool) error {
	return db.firestore.SaveEmailOTP(ctx, email, optIn)
}

// StageStartEpisodeByBreakGlass persists starts an emergency episode data
func (db DBService) StageStartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error {
	return db.firestore.StageStartEpisodeByBreakGlass(ctx, input)
}
