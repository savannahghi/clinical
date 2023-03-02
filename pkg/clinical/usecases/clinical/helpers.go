package clinical

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
)

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"text/template"

// 	"github.com/savannahghi/clinical/pkg/clinical/application/common"
// 	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
// 	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
// 	"github.com/savannahghi/clinical/pkg/clinical/domain"
// 	"github.com/savannahghi/enumutils"
// 	"github.com/savannahghi/scalarutils"
// 	"github.com/savannahghi/serverutils"
// 	log "github.com/sirupsen/logrus"
// )

// // TODO: remove receiver
// func validateSearchParams(params map[string]interface{}) (url.Values, error) {
// 	if params == nil {
// 		return nil, fmt.Errorf("can't search with nil params")
// 	}
// 	output := url.Values{}
// 	for k, v := range params {
// 		val, ok := v.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("the search/filter params should all be sent as strings")
// 		}
// 		output.Add(k, val)
// 	}
// 	return output, nil
// }

// // searchFilterHelper helps with the composition of FHIR REST search and filter requests.
// //
// // - the `resourceName` is a FHIR resource name e.g "Patient", "Appointment" etc
// // - the `path` is a resource sub-path e.g "_search". If there is no sub-path, send a blank string
// // - `params` should contain the filter parameters e.g
// //
// //    params := url.Values{}
// //    params.Add("_content", search)
// // TODO: remove receiver
// func (fh *UseCasesClinicalImpl) searchFilterHelper(
// 	ctx context.Context,
// 	resourceName string,
// 	path string, params url.Values,
// ) ([]map[string]interface{}, error) {
// 	// s.checkPreconditions()
// 	bs, err := fh.infrastructure.FHIR.POSTRequest(resourceName, path, params, nil)
// 	if err != nil {
// 		log.Errorf("unable to search: %v", err)
// 		return nil, fmt.Errorf("unable to search: %v", err)
// 	}
// 	respMap := make(map[string]interface{})
// 	err = json.Unmarshal(bs, &respMap)
// 	if err != nil {
// 		log.Errorf("%s could not be found with search params %v: %s", resourceName, params, err)
// 		return nil, fmt.Errorf(
// 			"%s could not be found with search params %v: %s", resourceName, params, err)
// 	}

// 	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
// 	for _, k := range mandatoryKeys {
// 		_, found := respMap[k]
// 		if !found {
// 			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
// 		}
// 	}
// 	resourceType, ok := respMap["resourceType"].(string)
// 	if !ok {
// 		return nil, fmt.Errorf("server error: the resourceType is not a string")
// 	}
// 	if resourceType != "Bundle" {
// 		return nil, fmt.Errorf(
// 			"server error: the resourceType value is not 'Bundle' as expected")
// 	}

// 	resultType, ok := respMap["type"].(string)
// 	if !ok {
// 		return nil, fmt.Errorf("server error: the search result type value is not a string")
// 	}
// 	if resultType != "searchset" {
// 		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
// 	}

// 	respEntries := respMap["entry"]
// 	if respEntries == nil {
// 		return []map[string]interface{}{}, nil
// 	}
// 	entries, ok := respEntries.([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf(
// 			"server error: entries is not a list of maps, it is: %T", respEntries)
// 	}

// 	results := []map[string]interface{}{}
// 	for _, en := range entries {
// 		entry, ok := en.(map[string]interface{})
// 		if !ok {
// 			return nil, fmt.Errorf(
// 				"server error: expected each entry to be map, they are %T instead", en)
// 		}
// 		expectedKeys := []string{"fullUrl", "resource", "search"}
// 		for _, k := range expectedKeys {
// 			_, found := entry[k]
// 			if !found {
// 				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
// 			}
// 		}

// 		resource, ok := entry["resource"].(map[string]interface{})
// 		if !ok {
// 			{
// 				return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
// 			}
// 		}
// 		results = append(results, resource)
// 	}
// 	return results, nil
// }

// TODO: remove receiver
func (c *UseCasesClinicalImpl) getTimelineEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, string, error) {
	episodePayload, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	episode := episodePayload.Resource
	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	if episode.Patient == nil || episode.Patient.Reference == nil {
		return nil, "", fmt.Errorf("the episode with ID %s has no patient reference", episodeID)
	}
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return nil, "", fmt.Errorf("the episode with ID %s is not active", episodeID)
	}
	if episode.Type == nil {
		return nil, "", fmt.Errorf("the episode with ID %s has a nil type", episodeID)
	}
	if len(episode.Type) != 1 {
		return nil, "", fmt.Errorf("expected the episode type to have just one entry")
	}
	accessLevel := episode.Type[0].Text
	if accessLevel != "FULL_ACCESS" && accessLevel != "PROFILE_AND_RECENT_VISITS_ACCESS" {
		return nil, "", fmt.Errorf("unknown episode access level: %s", accessLevel)
	}
	return episode, accessLevel, nil
}

// TODO: remove receiver
func (c *UseCasesClinicalImpl) getTimelineVisitSummaries(
	ctx context.Context,
	encounterSearchParams map[string]interface{},
	count int,
) ([]map[string]interface{}, error) {
	encounterConn, err := c.infrastructure.FHIR.SearchFHIREncounter(ctx, encounterSearchParams)
	if err != nil {
		return nil, fmt.Errorf("unable to search for encounters for the timeline: %w", err)
	}
	visitSummaries := []map[string]interface{}{}
	if encounterConn == nil || encounterConn.Edges == nil {
		return visitSummaries, nil
	}
	for _, edge := range encounterConn.Edges {
		if edge.Node == nil || edge.Node.ID == nil {
			continue
		}
		summary, err := c.VisitSummary(ctx, *edge.Node.ID, count)
		if err != nil {
			return nil, err
		}
		if summary != nil {
			visitSummaries = append(visitSummaries, summary)
		}
	}
	return visitSummaries, nil
}

// sendAlertToPatient to send notification to patient when break glass request is made
// TODO: move to engagement
func (c *UseCasesClinicalImpl) sendAlertToPatient(ctx context.Context, phoneNumber string, patientID string) error {
	patientPayload, err := c.FindPatientByID(ctx, patientID)
	if err != nil {
		return err
	}

	name := patientPayload.PatientRecord.Name[0].Given[0]
	if name == nil {
		return fmt.Errorf("nil patient name")
	}
	text := createAlertMessage(*name)

	type PayloadRequest struct {
		To      []string           `json:"to"`
		Message string             `json:"message"`
		Sender  enumutils.SenderID `json:"sender"`
	}

	requestPayload := PayloadRequest{
		To:      []string{phoneNumber},
		Message: text,
		Sender:  enumutils.SenderIDBewell,
	}

	resp, err := isc.MakeRequest(ctx, http.MethodPost, common.SendSMSEndpoint, requestPayload)
	if err != nil {
		return fmt.Errorf("unable to send alert to patient: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error status %s ", resp.Status)
	}

	return nil
}

// sendAlertToNextOfKin send an alert message to the patient's next of kin.
// TODO: move to engagement
func (c *UseCasesClinicalImpl) sendAlertToNextOfKin(ctx context.Context, patientID string) error {
	patientPayload, err := c.FindPatientByID(ctx, patientID)
	if err != nil {
		return err
	}
	patientContacts := patientPayload.PatientRecord.Contact
	patientName := patientPayload.PatientRecord.Name[0].Given[0]
	phone := domain.ContactPointSystemEnumPhone

	for _, patientRelation := range patientContacts {
		if patientRelation.Name == nil {
			continue
		}
		if len(patientRelation.Name.Given) == 0 {
			continue
		}
		for _, codeableConcept := range patientRelation.Relationship {
			for _, coding := range codeableConcept.Coding {
				if coding.Code == "N" {
					// this is the next of kin
					phoneNextOfKin := patientRelation.Telecom
					for _, number := range phoneNextOfKin {
						if number == nil {
							continue
						}
						if number.Value == nil {
							continue
						}

						numberSystem := *number.System

						if &numberSystem == &phone {
							text := createNextOfKinAlertMessage(*patientRelation.Name.Given[0], *patientName)

							type PayloadRequest struct {
								To      []string           `json:"to"`
								Message string             `json:"message"`
								Sender  enumutils.SenderID `json:"sender"`
							}

							requestPayload := PayloadRequest{
								To:      []string{*number.Value},
								Message: text,
								Sender:  enumutils.SenderIDBewell,
							}

							resp, err := isc.MakeRequest(ctx, http.MethodPost, common.SendSMSEndpoint, requestPayload)
							if err != nil {
								return fmt.Errorf("unable to send alert to next of kin: %w", err)
							}

							if resp.StatusCode != http.StatusOK {
								return fmt.Errorf("got error status %s from email service", resp.Status)
							}
							return nil
						}
					}
					break
				}
			}

		}
	}
	err = fmt.Errorf("failed to send message")
	return err
}

// sendAlertToAdmin send email to admin notifying them of the access.
// TODO: move to engagement
func (c *UseCasesClinicalImpl) sendAlertToAdmin(ctx context.Context, patientName string, patientContact string) error {
	adminEmail, err := serverutils.GetEnvVar(SavannahAdminEmail)
	if err != nil {
		return err
	}

	var writer bytes.Buffer
	t := template.Must(template.New("profile").Parse(utils.AdminEmailMessage))
	_ = t.Execute(&writer, struct {
		Name   string
		Number string
	}{
		Name:   patientName,
		Number: patientContact,
	})
	subject := "Breaking Glass Access notice"

	body := map[string]interface{}{
		"to":      []string{adminEmail},
		"text":    writer.String(),
		"subject": subject,
	}

	resp, err := isc.MakeRequest(ctx, http.MethodPost, common.SendEmailEndpoint, body)
	if err != nil {
		return fmt.Errorf("unable to send Alert to admin email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("unable to send Alert to admin email : %v, with status code %v", err, resp.StatusCode)
		return fmt.Errorf("got error status %s from email service", resp.Status)
	}

	return nil
}

// createAlertMessage Create a nice message to be sent.
func createAlertMessage(names string) string {
	text := fmt.Sprintf(
		"Dear %s. Your health record has been accessed for an emergency. "+
			"If you are not aware of the circumstances of this, please call %s",
		names,
		common.CallCenterNumber,
	)
	return text
}

// createNextOfKinAlertMessage creates a message to be sent to the next of kin
func createNextOfKinAlertMessage(names, patientName string) string {
	text := fmt.Sprintf(
		"Dear %s. The health record for %s has been accessed for an emergency. "+
			"If you are not aware of the circumstances of this, please call %s",
		names,
		patientName,
		common.CallCenterNumber,
	)
	return text
}
