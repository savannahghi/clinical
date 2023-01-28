package pubsubmessaging_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	pubsubmessaging "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

func InitializeTestPubSub(t *testing.T) (*pubsubmessaging.ServicePubSubMessaging, error) {
	fc := &firebasetools.FirebaseClient{}
	ctx := context.Background()
	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}

	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	// Initialize base (common) extension
	baseExtension := extensions.NewBaseExtensionImpl(fc)
	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl)
	usecases := usecases.NewUsecasesInteractor(infrastructure)
	oclUseCase := usecases.UseCasesOCL
	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		baseExtension,
		infrastructure,
		usecases,
		oclUseCase,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pubsub messaging service: %w", err)
	}
	return pubSub, nil
}

func TestServicePubSubMessaging_AddPubSubNamespace(t *testing.T) {
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	topicName := pubsubmessaging.TestTopicName
	environment := serverutils.GetRunningEnvironment()

	type args struct {
		topicName   string
		serviceName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case -> Correct pubsub namespace",
			args: args{
				topicName:   topicName,
				serviceName: pubsubmessaging.ClinicalServiceName,
			},
			want: fmt.Sprintf("%s-%s-%s-%s",
				pubsubmessaging.ClinicalServiceName,
				topicName,
				environment,
				pubsubmessaging.TopicVersion,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.AddPubSubNamespace(tt.args.topicName, tt.args.serviceName)
			if got != tt.want {
				t.Errorf("ServicePubSubMessaging.AddPubSubNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServicePubSubMessaging_PublishToPubsub(t *testing.T) {
	ctx := context.Background()
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	topic := ps.AddPubSubNamespace(pubsubmessaging.TestTopicName, pubsubmessaging.ClinicalServiceName)
	// Create the test topic
	topics := ps.TopicIDs()
	topics = append(topics, topic)

	err = ps.EnsureTopicsExist(ctx, topics)
	if err != nil {
		t.Errorf("failed to create test topic")
		return
	}

	payload := map[string]interface{}{
		"name": "Test PubsubPayload",
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		topicID string
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad Case -> Fail to publish to pubsub - nil payload",
			args: args{
				ctx:     ctx,
				topicID: topic,
				payload: nil,
			},
			wantErr: true,
		},
		{
			name: "Sad Case -> Fail to publish to pubsub - unknown topic",
			args: args{
				ctx:     ctx,
				topicID: "invalid",
				payload: marshalled,
			},
			wantErr: true,
		},
		{
			name: "Happy Case-> Publish to pubsub",
			args: args{
				ctx:     ctx,
				topicID: topic,
				payload: marshalled,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.PublishToPubsub(tt.args.ctx, tt.args.topicID, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.PublishToPubsub() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_EnsureTopicsExist(t *testing.T) {
	ctx := context.Background()
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	type args struct {
		ctx      context.Context
		topicIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Create a topic",
			args: args{
				ctx:      ctx,
				topicIDs: []string{pubsubmessaging.TestTopicName},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.EnsureTopicsExist(tt.args.ctx, tt.args.topicIDs); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.EnsureTopicsExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_SubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}
	topic := ps.AddPubSubNamespace(pubsubmessaging.TestTopicName, pubsubmessaging.ClinicalServiceName)
	// Create the test topic
	topics := ps.TopicIDs()
	topics = append(topics, topic)

	err = ps.EnsureTopicsExist(ctx, topics)
	if err != nil {
		t.Errorf("failed to create test topic")
		return
	}

	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "Happy Case -> return subscriptionIDs",
			want: map[string]string{
				topic: fmt.Sprintf("%s-default-subscription", topics),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.SubscriptionIDs()
			if got == nil {
				t.Errorf("ServicePubSubMessaging.SubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
