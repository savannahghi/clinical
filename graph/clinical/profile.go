package clinical

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bearbin/go-age"
	"github.com/savannahghi/scalarutils"
)

// RenderOfficialName returns a patient's official name in Markdown format
func (p FHIRPatient) RenderOfficialName() scalarutils.Markdown {
	defaultName := scalarutils.Markdown("UNKNOWN NAME")
	if p.Name == nil {
		return defaultName
	}
	// prefer official names
	for _, name := range p.Name {
		if name != nil && name.Use == HumanNameUseEnumOfficial {
			return scalarutils.Markdown(name.Text)
		}
	}
	// fall back to usual names
	for _, name := range p.Name {
		if name != nil && name.Use == HumanNameUseEnumUsual {
			return scalarutils.Markdown(name.Text)
		}
	}
	return defaultName
}

// RenderIDDocuments returns a patient's ID documents in Markdown format
func (p FHIRPatient) RenderIDDocuments() scalarutils.Markdown {
	defaultID := scalarutils.Markdown("No Identification documents found")
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
		return scalarutils.Markdown(ids)
	}
	return defaultID
}

// RenderAge returns the patient's age, rendered in a humanized way
func (p FHIRPatient) RenderAge() scalarutils.Markdown {
	ageStr := "UNKNOWN AGE"
	if p.BirthDate != nil {
		ageStr = strconv.Itoa(age.Age(p.BirthDate.AsTime())) + " yrs"
	}
	return scalarutils.Markdown(fmt.Sprintf("Age: %s", ageStr))
}

// RenderGender returns the patient's age and gender, rendered in a humanized way
func (p FHIRPatient) RenderGender() scalarutils.Markdown {
	gender := "UNKNOWN GENDER"
	if p.Gender != nil {
		gender = p.Gender.String()
	}

	return scalarutils.Markdown(fmt.Sprintf("Gender: %s", gender))
}

// RenderProblems returns the patient's problems
func (p FHIRPatient) RenderProblems() scalarutils.Markdown {
	defaultProblem := scalarutils.Markdown("Problems: No known problems")
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
		return scalarutils.Markdown(patientProblems)
	}

	return scalarutils.Markdown(defaultProblem)
}

func trimString(inp string, maxLength int) string {
	if len(inp) <= maxLength {
		return inp
	}
	return inp[:maxLength-3] + "..."
}
