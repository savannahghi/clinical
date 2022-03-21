package firebase

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
)

const (
	// EmailOptInCollectionName ...
	EmailOptInCollectionName = "email_opt_ins"
	// BreakGlassCollectionName ...
	BreakGlassCollectionName = "break_glass"
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	FirestoreClient FirestoreClientExtension
	FirebaseClient  FBClientExtension
}

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(
	firestoreClient FirestoreClientExtension,
	firebaseClient FBClientExtension,
) *Repository {
	return &Repository{
		FirestoreClient: firestoreClient,
		FirebaseClient:  firebaseClient,
	}
}

// GetEmailOptInCollectionName ...
func (fr Repository) GetEmailOptInCollectionName() string {

	suffixed := firebasetools.SuffixCollection(EmailOptInCollectionName)
	return suffixed
}

func (fr Repository) getBreakGlassCollectionName() string {
	suffixed := firebasetools.SuffixCollection(BreakGlassCollectionName)
	return suffixed
}

// StageStartEpisodeByBreakGlass persists starts an emergency episode data
func (fr Repository) StageStartEpisodeByBreakGlass(
	ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error {
	command := &CreateCommand{
		CollectionName: fr.getBreakGlassCollectionName(),
		Data:           input,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return fmt.Errorf("unable to stage start episode by break glass data: %w", err)
	}
	return nil
}

// SaveEmailOTP  persist the data of the newly created OTP to a datastore
func (fr Repository) SaveEmailOTP(
	ctx context.Context,
	email string, optIn bool) error {
	if email == "" {
		return fmt.Errorf("the email cannot be nil")
	}
	if optIn {
		data := domain.EmailOptIn{
			Email:   email,
			OptedIn: optIn,
		}
		command := &CreateCommand{
			CollectionName: fr.GetEmailOptInCollectionName(),
			Data:           data,
		}
		_, err := fr.FirestoreClient.Create(ctx, command)
		if err != nil {
			return fmt.Errorf("unable to save email opt in: %w", err)
		}
	}
	return nil
}
