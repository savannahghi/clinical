package fhir

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/scalarutils"
)

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
