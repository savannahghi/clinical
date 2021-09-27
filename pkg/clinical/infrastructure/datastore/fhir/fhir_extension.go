package fhir

// constants used to configure the Google Cloud Healthcare API
const (
	DatasetLocation = "europe-west4"
	baseFHIRURL     = "https://healthcare.googleapis.com/v1"
)

// // DatasetExtension represents the methods we need from Healthcare Dataset and FHIR service
// type DatasetExtension interface {
// 	GetDataset() (*healthcare.Dataset, error)
// 	CreateDataset() (*healthcare.Operation, error)
// 	GetFHIRStore() (*healthcare.FhirStore, error)
// 	CreateFHIRStore() (*healthcare.FhirStore, error)
// }

// // DatasetExtensionImpl ...
// type DatasetExtensionImpl struct {
// 	healthcareService                           *healthcare.Service
// 	projectID, location, datasetID, fhirStoreID string
// }

// NewDatasetExtension initializes a new Dataset extension
// func NewDatasetExtension() DatasetExtension {
// 	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
// 	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
// 	dataset := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
// 	fhirStore := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
// 	ctx := context.Background()
// 	hsv, err := healthcare.NewService(ctx)
// 	if err != nil {
// 		log.Panicf("unable to initialize new Google Cloud Healthcare Service: %s", err)
// 	}
// 	return &DatasetExtensionImpl{
// 		healthcareService: hsv,
// 		projectID:         project,
// 		location:          DatasetLocation,
// 		datasetID:         dataset,
// 		fhirStoreID:       fhirStore,
// 	}
// }

// // CreateDataset creates a dataset and returns it's name
// func (d DatasetExtensionImpl) CreateDataset() (*healthcare.Operation, error) {
// 	d.checkPreconditions()
// 	datasetsService := d.healthcareService.Projects.Locations.Datasets
// 	parent := fmt.Sprintf("projects/%s/locations/%s", d.projectID, d.location)
// 	resp, err := datasetsService.Create(parent, &healthcare.Dataset{}).DatasetId(d.datasetID).Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("create Data Set: %v", err)
// 	}
// 	return resp, nil
// }

// // GetDataset gets a dataset.
// func (d DatasetExtensionImpl) GetDataset() (*healthcare.Dataset, error) {
// 	d.checkPreconditions()
// 	datasetsService := d.healthcareService.Projects.Locations.Datasets
// 	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", d.projectID, d.location, d.datasetID)
// 	resp, err := datasetsService.Get(name).Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("get Data Set: %v", err)
// 	}
// 	return resp, nil
// }

// // CreateFHIRStore creates an FHIR store.
// func (d DatasetExtensionImpl) CreateFHIRStore() (*healthcare.FhirStore, error) {
// 	d.checkPreconditions()
// 	storesService := d.healthcareService.Projects.Locations.Datasets.FhirStores
// 	store := &healthcare.FhirStore{
// 		DisableReferentialIntegrity: false,
// 		DisableResourceVersioning:   false,
// 		EnableUpdateCreate:          true,
// 		Version:                     "R4",
// 	}
// 	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", d.projectID, d.location, d.datasetID)
// 	resp, err := storesService.Create(parent, store).FhirStoreId(d.fhirStoreID).Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("create FHIR Store: %v", err)
// 	}
// 	return resp, nil
// }

// // GetFHIRStore gets an FHIR store.
// func (d DatasetExtensionImpl) GetFHIRStore() (*healthcare.FhirStore, error) {
// 	d.checkPreconditions()
// 	storesService := d.healthcareService.Projects.Locations.Datasets.FhirStores
// 	name := fmt.Sprintf(
// 		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
// 		d.projectID, d.location, d.datasetID, d.fhirStoreID)
// 	store, err := storesService.Get(name).Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("get FHIR Store: %v", err)
// 	}
// 	return store, nil
// }

// func (d DatasetExtensionImpl) checkPreconditions() {
// 	if d.healthcareService == nil {
// 		log.Panicf("cloudhealth.Service *healthcare.Service is nil")
// 	}
// }
