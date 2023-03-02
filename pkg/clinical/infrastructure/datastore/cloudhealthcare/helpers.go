package fhir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/labstack/gommon/log"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/scalarutils"
)

func validateSearchParams(params map[string]interface{}) (url.Values, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	output := url.Values{}
	for k, v := range params {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("the search/filter params should all be sent as strings")
		}
		output.Add(k, val)
	}
	return output, nil
}

// searchFilterHelper helps with the composition of FHIR REST search and filter requests.
//
// - the `resourceName` is a FHIR resource name e.g "Patient", "Appointment" etc
// - the `path` is a resource sub-path e.g "_search". If there is no sub-path, send a blank string
// - `params` should contain the filter parameters e.g
//
//	params := url.Values{}
//	params.Add("_content", search)
//
// TODO: remove receiver
func (fh *StoreImpl) searchFilterHelper(
	ctx context.Context,
	resourceName string,
	path string, params url.Values,
) ([]map[string]interface{}, error) {
	// s.checkPreconditions()
	bs, err := fh.Dataset.POSTRequest(resourceName, path, params, nil)
	if err != nil {
		log.Errorf("unable to search: %v", err)
		return nil, fmt.Errorf("unable to search: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		log.Errorf("%s could not be found with search params %v: %s", resourceName, params, err)
		return nil, fmt.Errorf(
			"%s could not be found with search params %v: %s", resourceName, params, err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		return nil, fmt.Errorf(
			"server error: the resourceType value is not 'Bundle' as expected")
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the search result type value is not a string")
	}
	if resultType != "searchset" {
		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return []map[string]interface{}{}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	results := []map[string]interface{}{}
	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
			}
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			{
				return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
			}
		}
		results = append(results, resource)
	}
	return results, nil
}

func birthdateMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	parsedDate := helpers.ParseDate(resourceCopy["birthDate"].(string))

	dateMap := make(map[string]interface{})

	dateMap["year"] = parsedDate.Year()
	dateMap["month"] = parsedDate.Month()
	dateMap["day"] = parsedDate.Day()

	resourceCopy["birthDate"] = dateMap

	return resourceCopy
}

func periodMapper(period map[string]interface{}) map[string]interface{} {

	periodCopy := period

	parsedStartDate := helpers.ParseDate(periodCopy["start"].(string))

	periodCopy["start"] = scalarutils.DateTime(parsedStartDate.Format(timeFormatStr))

	parsedEndDate := helpers.ParseDate(periodCopy["end"].(string))

	periodCopy["end"] = scalarutils.DateTime(parsedEndDate.Format(timeFormatStr))

	return periodCopy
}

func identifierMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	if _, ok := resource["identifier"]; ok {

		newIdentifiers := []map[string]interface{}{}

		for _, identifier := range resource["identifier"].([]interface{}) {

			identifier := identifier.(map[string]interface{})

			if _, ok := identifier["period"]; ok {

				period := identifier["period"].(map[string]interface{})
				newPeriod := periodMapper(period)

				identifier["period"] = newPeriod
			}

			newIdentifiers = append(newIdentifiers, identifier)
		}

		resourceCopy["identifier"] = newIdentifiers
	}

	return resourceCopy
}

func nameMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newNames := []map[string]interface{}{}

	if _, ok := resource["name"]; ok {

		for _, name := range resource["name"].([]interface{}) {

			name := name.(map[string]interface{})

			if _, ok := name["period"]; ok {

				period := name["period"].(map[string]interface{})
				newPeriod := periodMapper(period)

				name["period"] = newPeriod
			}

			newNames = append(newNames, name)
		}

	}

	resourceCopy["name"] = newNames

	return resourceCopy
}

func telecomMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newTelecoms := []map[string]interface{}{}

	if _, ok := resource["telecom"]; ok {

		for _, telecom := range resource["telecom"].([]interface{}) {

			telecom := telecom.(map[string]interface{})

			if _, ok := telecom["period"]; ok {

				period := telecom["period"].(map[string]interface{})
				newPeriod := periodMapper(period)

				telecom["period"] = newPeriod
			}

			newTelecoms = append(newTelecoms, telecom)
		}

	}

	resourceCopy["telecom"] = newTelecoms

	return resourceCopy
}

func addressMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newAddresses := []map[string]interface{}{}

	if _, ok := resource["address"]; ok {

		for _, address := range resource["address"].([]interface{}) {

			address := address.(map[string]interface{})

			if _, ok := address["period"]; ok {

				period := address["period"].(map[string]interface{})
				newPeriod := periodMapper(period)

				address["period"] = newPeriod
			}

			newAddresses = append(newAddresses, address)
		}
	}

	resourceCopy["address"] = newAddresses

	return resourceCopy
}

func photoMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newPhotos := []map[string]interface{}{}

	if _, ok := resource["photo"]; ok {

		for _, photo := range resource["photo"].([]interface{}) {

			photo := photo.(map[string]interface{})

			parsedDate := helpers.ParseDate(photo["creation"].(string))

			photo["creation"] = scalarutils.DateTime(parsedDate.Format(timeFormatStr))

			newPhotos = append(newPhotos, photo)
		}
	}

	resourceCopy["photo"] = newPhotos

	return resourceCopy
}

func contactMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newContacts := []map[string]interface{}{}

	if _, ok := resource["contact"]; ok {

		for _, contact := range resource["contact"].([]interface{}) {

			contact := contact.(map[string]interface{})

			if _, ok := contact["name"]; ok {

				name := contact["name"].(map[string]interface{})
				if _, ok := name["period"]; ok {

					period := name["period"].(map[string]interface{})
					newPeriod := periodMapper(period)

					name["period"] = newPeriod
				}

				contact["name"] = name
			}

			if _, ok := contact["telecom"]; ok {

				newTelecoms := []map[string]interface{}{}

				for _, telecom := range contact["telecom"].([]interface{}) {

					telecom := telecom.(map[string]interface{})

					if _, ok := telecom["period"]; ok {

						period := telecom["period"].(map[string]interface{})
						newPeriod := periodMapper(period)

						telecom["period"] = newPeriod
					}

					newTelecoms = append(newTelecoms, telecom)
				}

				contact["telecom"] = newTelecoms
			}

			if _, ok := contact["period"]; ok {

				period := contact["period"].(map[string]interface{})
				newPeriod := periodMapper(period)

				contact["period"] = newPeriod
			}

			newContacts = append(newContacts, contact)
		}
	}

	resourceCopy["contact"] = newContacts

	return resourceCopy
}
