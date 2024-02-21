package advantage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/serverutils"
)

var (
	AdvantageBaseURL = serverutils.MustGetEnvVar("ADVANTAGE_BASE_URL")
	segmentationPath = "/api/segments/segment/clinical/"
)

// AuthUtilsLib holds the method defined in authutils library
type AuthUtilsLib interface {
	Authenticate() (*authutils.OAUTHResponse, error)
}

// AdvantageService represents methods that can be used to communicate with the advantage server
type AdvantageService interface {
	SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error
}

// ServiceAdvantageImpl represents advantage server's implementations
type ServiceAdvantageImpl struct {
	client AuthUtilsLib
}

// NewServiceAdvantage is the advantage server's service constructor
func NewServiceAdvantage(authUtils AuthUtilsLib) *ServiceAdvantageImpl {
	return &ServiceAdvantageImpl{
		client: authUtils,
	}
}

// SegmentPatient is used to create segmentation information in advantage
func (s *ServiceAdvantageImpl) SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error {
	url := fmt.Sprintf("%s/%s", AdvantageBaseURL, segmentationPath)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	token, err := s.client.Authenticate()
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	httpClient := &http.Client{Timeout: time.Second * 30}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
