package firebase

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
)

const (
	// EmailOptInCollectionName ...
	EmailOptInCollectionName = "email_opt_ins"
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

// ValidateEmail returns an error if the supplied string does not have a
// valid format or resolvable host
func (fr Repository) ValidateEmail(
	ctx context.Context,
	email string, optIn bool) error {

	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
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
