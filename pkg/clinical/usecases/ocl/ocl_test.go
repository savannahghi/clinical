package ocl_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

var (
	testUsecaseInteractor interactor.Usecases
	testInfrastructure    infrastructure.Infrastructure
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")

	log.Printf("Running tests ...")

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}

	testInfrastructure = infra

	svc, err := InitializeTestService(ctx)
	if err != nil {
		log.Printf("failed to initialize test service: %v", err)
	}
	testUsecaseInteractor = svc

	code := m.Run()

	log.Printf("Tearing tests down ...")

	os.Exit(code)
}

func InitializeTestService(ctx context.Context) (interactor.Usecases, error) {
	fc := &firebasetools.FirebaseClient{}
	baseExtension := extensions.NewBaseExtensionImpl(fc)
	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl)

	usecases := interactor.NewUsecasesInteractor(
		infrastructure,
	)
	return usecases, nil
}

func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Infrastructure, error) {
	fc := &firebasetools.FirebaseClient{}
	baseExtension := extensions.NewBaseExtensionImpl(fc)
	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)

	ocl := openconceptlab.NewServiceOCL()
	return infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl), nil
}

func TestUseCasesImpl_ListConcepts(t *testing.T) {
	svc := testUsecaseInteractor

	type args struct {
		ctx                    context.Context
		org                    string
		source                 string
		verbose                bool
		q                      *string
		sortAsc                *string
		sortDesc               *string
		conceptClass           *string
		dataType               *string
		locale                 *string
		includeRetired         *bool
		includeMappings        *bool
		includeInverseMappings *bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:     context.Background(),
				org:     "CIEL",
				source:  "CIEL",
				verbose: true,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:     context.Background(),
				org:     "this is unknown org",
				source:  "CIEL",
				verbose: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.ListConcepts(tt.args.ctx, tt.args.org, tt.args.source, tt.args.verbose, tt.args.q, tt.args.sortAsc, tt.args.sortDesc, tt.args.conceptClass, tt.args.dataType, tt.args.locale, tt.args.includeRetired, tt.args.includeMappings, tt.args.includeInverseMappings)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListConcepts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestUseCasesImpl_GetConcept(t *testing.T) {
	svc := testUsecaseInteractor
	type args struct {
		ctx                    context.Context
		org                    string
		source                 string
		concept                string
		includeMappings        bool
		includeInverseMappings bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid_who_icd_concept",
			args: args{
				ctx:                    context.Background(),
				org:                    "CIEL",
				source:                 "CIEL",
				concept:                "106",
				includeMappings:        true,
				includeInverseMappings: false,
			},
			wantErr: false,
		},
		{
			name: "sad case -- nil concept in the payload",
			args: args{
				ctx:                    context.Background(),
				org:                    "CIEL",
				source:                 "CIEL",
				includeMappings:        true,
				includeInverseMappings: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ocl := svc
			got, err := ocl.GetConcept(tt.args.ctx, tt.args.org, tt.args.source, tt.args.concept, tt.args.includeMappings, tt.args.includeInverseMappings)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesImpl.GetConcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Contains(t, got, "display_name")
			}
		})
	}
}
