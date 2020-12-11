package clinical

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bearbin/go-age"
	"gitlab.slade360emr.com/go/base"
)

// RenderOfficialName returns a patient's official name in Markdown format
func (p FHIRPatient) RenderOfficialName() base.Markdown {
	defaultName := base.Markdown("UNKNOWN NAME")
	if p.Name == nil {
		return defaultName
	}
	// prefer official names
	for _, name := range p.Name {
		if name != nil && name.Use == HumanNameUseEnumOfficial {
			return base.Markdown(name.Text)
		}
	}
	// fall back to usual names
	for _, name := range p.Name {
		if name != nil && name.Use == HumanNameUseEnumUsual {
			return base.Markdown(name.Text)
		}
	}
	return defaultName
}

// RenderIDDocuments returns a patient's ID documents in Markdown format
func (p FHIRPatient) RenderIDDocuments() base.Markdown {
	defaultID := base.Markdown("No Identification documents found")
	identifiers := []string{}
	if p.Identifier != nil {
		for _, identifier := range p.Identifier {
			if identifier == nil {
				continue
			}
			identif := fmt.Sprintf("%s (%s)", identifier.Value, identifier.Use.String())
			identifiers = append(identifiers, identif)
		}
	}
	if len(identifiers) > 0 {
		ids := strings.Join(identifiers, " , ")
		return base.Markdown(ids)
	}
	return defaultID
}

// RenderAge returns the patient's age, rendered in a humanized way
func (p FHIRPatient) RenderAge() base.Markdown {
	ageStr := "UNKNOWN AGE"
	if p.BirthDate != nil {
		ageStr = strconv.Itoa(age.Age(p.BirthDate.AsTime())) + " yrs"
	}
	return base.Markdown(fmt.Sprintf("Age: %s", ageStr))
}

// RenderGender returns the patient's age and gender, rendered in a humanized way
func (p FHIRPatient) RenderGender() base.Markdown {
	gender := "UNKNOWN GENDER"
	if p.Gender != nil {
		gender = p.Gender.String()
	}

	return base.Markdown(fmt.Sprintf("Gender: %s", gender))
}

// RenderProblems returns the patient's problems
func (p FHIRPatient) RenderProblems() base.Markdown {
	defaultProblem := base.Markdown("Problems: No known problems")
	if p.ID == nil {
		return defaultProblem
	}
	clinicalService := NewService()
	ctx := context.Background()
	problemSummary, err := clinicalService.ProblemSummary(ctx, *p.ID)
	if err != nil {
		return defaultProblem
	}

	problems := ""
	if len(problemSummary) > 0 {
		problems = strings.Join(problemSummary, ",")
	}

	if len(problems) > 0 {
		patientProblems := fmt.Sprintf("Problems: %s", problems)
		return base.Markdown(patientProblems)
	}

	return base.Markdown(defaultProblem)
}

// RenderAllergies returns the patient's allergies
func (p FHIRPatient) RenderAllergies() base.Markdown {
	defaultAllergies := base.Markdown("Allergies: No known allergies")
	if p.ID == nil {
		return defaultAllergies
	}
	clinicalService := NewService()
	ctx := context.Background()
	allergySummary, err := clinicalService.AllergySummary(ctx, *p.ID)
	if err != nil {
		return defaultAllergies
	}

	allergies := ""
	if len(allergySummary) > 0 {
		allergies = strings.Join(allergySummary, ",")
	}

	if len(allergies) > 0 {
		patientAllergies := fmt.Sprintf("Allergies: %s", allergies)
		return base.Markdown(patientAllergies)
	}
	return base.Markdown(defaultAllergies)
}

func trimString(inp string, maxLength int) string {
	if len(inp) <= maxLength {
		return inp
	}
	return inp[:maxLength-3] + "..."
}
