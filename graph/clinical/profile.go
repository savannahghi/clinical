package clinical

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bearbin/go-age"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/base"
)

// RequestUSSDLastVisit returns details of the patient's last visit
func (s Service) RequestUSSDLastVisit(
	ctx context.Context, input USSDClinicalRequest) (*USSDClinicalResponse, error) {
	s.checkPreconditions()

	patient, err := s.lookupUSSDSessionPatient(ctx, input)
	if err != nil {
		return nil, err
	}

	serverDomain, err := base.GetEnvVar(base.ServerPublicDomainEnvironmentVariableName)
	if err != nil {
		return nil, err
	}

	generatedURL, err := GeneratePatientLink(*patient.ID)
	if err != nil {
		return nil, err
	}

	longURL := fmt.Sprintf("%s/visits/%s", serverDomain, generatedURL.OpaqueID)
	shortURL, err := base.ShortenLink(ctx, longURL)
	if err != nil {
		return nil, fmt.Errorf("unable to shorten patient profile URL: %v", err)
	}

	name := patient.RenderOfficialName()
	text := fmt.Sprintf(
		"Dear %s, please access your patient record at %s",
		name, shortURL,
	)
	summary := fmt.Sprintf(
		"%s\n\nPlease access your visit summary at %s", name, shortURL)

	return &USSDClinicalResponse{
		ShortLink: shortURL,
		Summary:   summary,
		Text:      text,
	}, nil
}

// RequestUSSDFullHistory returns links to the patient's full history
func (s Service) RequestUSSDFullHistory(
	ctx context.Context, input USSDClinicalRequest) (*USSDClinicalResponse, error) {
	s.checkPreconditions()

	patient, err := s.lookupUSSDSessionPatient(ctx, input)
	if err != nil {
		return nil, err
	}

	serverDomain, err := base.GetEnvVar(base.ServerPublicDomainEnvironmentVariableName)
	if err != nil {
		return nil, err
	}

	generatedURL, err := GeneratePatientLink(*patient.ID)
	if err != nil {
		return nil, err
	}

	longURL := fmt.Sprintf("%s/charts/%s", serverDomain, generatedURL.OpaqueID)
	shortURL, err := base.ShortenLink(ctx, longURL)
	if err != nil {
		return nil, fmt.Errorf("unable to shorted patient profile URL: %v", err)
	}

	name := patient.RenderOfficialName()
	text := fmt.Sprintf(
		"Dear %s, please access your last visit's details at %s",
		name, shortURL,
	)

	summary := fmt.Sprintf(
		"%s\n\nPlease access your clinical history at %s", name, shortURL)

	return &USSDClinicalResponse{
		ShortLink: shortURL,
		Summary:   summary,
		Text:      text,
	}, nil
}

// RequestUSSDPatientProfile returns details of the patient's profile in a form
// that is suitable for use by a USSD gateway
func (s Service) RequestUSSDPatientProfile(
	ctx context.Context, input USSDClinicalRequest) (*USSDClinicalResponse, error) {
	s.checkPreconditions()

	patient, err := s.lookupUSSDSessionPatient(ctx, input)
	if err != nil {
		return nil, err
	}

	if patient.Active == nil || !*patient.Active {
		return nil, fmt.Errorf("inactive patient / nil active")
	}

	if patient.ID == nil {
		return nil, fmt.Errorf("nil patient ID for patient fetched using USSD input %#v", input)
	}

	serverDomain, err := base.GetEnvVar(base.ServerPublicDomainEnvironmentVariableName)
	if err != nil {
		return nil, err
	}

	generatedURL, err := GeneratePatientLink(*patient.ID)
	if err != nil {
		return nil, err
	}

	longURL := fmt.Sprintf("%s/profiles/%s", serverDomain, generatedURL.OpaqueID)
	shortURL, err := base.ShortenLink(ctx, longURL)
	if err != nil {
		return nil, fmt.Errorf("unable to shorted patient profile URL: %v", err)
	}

	name := patient.RenderOfficialName()
	age := patient.RenderAge()
	problems := patient.RenderProblems()
	allergies := patient.RenderAllergies()
	text := fmt.Sprintf(
		"Dear %s, please access your health profile at %s",
		name, shortURL,
	)

	summary := fmt.Sprintf(
		"%s (%s)\n\nProblems and allergies:\n%s,%s\nSee more at %s",
		name, age, problems, allergies, shortURL,
	)
	return &USSDClinicalResponse{
		ShortLink: shortURL,
		Summary:   summary,
		Text:      text,
	}, nil
}

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

// RenderAddresses renders the patient's postal and physical addresses
func (p FHIRPatient) RenderAddresses() base.Markdown {
	defaultAddress := base.Markdown("")
	if p.Address == nil {
		return defaultAddress
	}
	addresses := []string{}
	for _, addr := range p.Address {
		if addr == nil {
			continue
		}
		addresses = append(addresses, addr.Text)
	}
	return base.Markdown("\n" + strings.Join(addresses, "\n"))
}

// RenderMaritalStatus renders the patient's marital status
func (p FHIRPatient) RenderMaritalStatus() base.Markdown {
	defaultMaritalStatus := base.Markdown("UNKNOWN MARITAL STATUS")
	if p.MaritalStatus == nil {
		return base.Markdown(p.MaritalStatus.Text)
	}
	return defaultMaritalStatus
}

// RenderLanguages renders the patient's languages as Markdown
func (p FHIRPatient) RenderLanguages() base.Markdown {
	defaultLanguages := base.Markdown("-")
	if p.Communication == nil {
		return defaultLanguages
	}
	languages := []string{}
	for _, comm := range p.Communication {
		if comm == nil {
			continue
		}
		if comm.Language == nil {
			continue
		}
		if comm.Preferred != nil && *comm.Preferred {
			languages = append(languages, fmt.Sprintf("%s (preferred)", comm.Language.Text))
		} else {
			languages = append(languages, comm.Language.Text)
		}
	}
	return base.Markdown(strings.Join(languages, ", "))
}

// RenderVisitSummary renders the patient's visit summary
func (p FHIRPatient) RenderVisitSummary(ctx context.Context, clinicalService *Service) map[string]interface{} {
	encounterSearchParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%s", *p.ID),
		"_sort":   "date",
		"count":   1,
	}
	defaultVisitSummary := make(map[string]interface{})
	encounterConnection, err := clinicalService.SearchFHIREncounter(ctx, encounterSearchParams)
	if err != nil {
		return defaultVisitSummary
	}
	if len(encounterConnection.Edges) == 0 {
		return defaultVisitSummary
	}
	encounter := encounterConnection.Edges[0].Node
	visitSummary, err := clinicalService.VisitSummary(ctx, *encounter.ID)
	if err != nil {
		log.Printf("Unable to render visit summary for encounter %#v", *encounter)
		return defaultVisitSummary
	}

	return visitSummary
}

// RenderFullHistory renders the patient's timeline
func (p FHIRPatient) RenderFullHistory(ctx context.Context, clinicalService *Service) []map[string]interface{} {
	episodeSearchParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%s", *p.ID),
		"_sort":   "date",
		"count":   1,
	}
	defaultEpisodeSummary := make([]map[string]interface{}, 0)
	episodeConnection, err := clinicalService.SearchFHIREpisodeOfCare(ctx, episodeSearchParams)
	if err != nil {
		return defaultEpisodeSummary
	}
	if len(episodeConnection.Edges) == 0 {
		return defaultEpisodeSummary
	}
	episode := episodeConnection.Edges[0].Node
	timeline, err := clinicalService.PatientTimeline(ctx, *episode.ID)
	if err != nil {
		log.Printf("Unable to render timeline for episode %#v", *episode)
		return defaultEpisodeSummary
	}
	return timeline
}

// PatientProfileHandlerFunc returns a function that renders a patient's profile page
func PatientProfileHandlerFunc(clinicalService *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		ctx := context.Background()
		patientID, err := GetPatientID(ctx, id)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData)
		}
		payload, err := clinicalService.GetFHIRPatient(ctx, patientID)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData) // we are already trying to handle an error
		}
		patient := payload.Resource
		name := patient.RenderOfficialName()
		age := patient.RenderAge()
		gender := patient.RenderGender()
		problems := patient.RenderProblems()
		allergies := patient.RenderAllergies()
		idDocs := patient.RenderIDDocuments()
		addresses := patient.RenderAddresses()
		maritalStatus := patient.RenderMaritalStatus()
		languages := patient.RenderLanguages()

		templateData := struct {
			Name          base.Markdown
			Age           base.Markdown
			Gender        base.Markdown
			Problems      base.Markdown
			Allergies     base.Markdown
			IDDocs        base.Markdown
			Addresses     base.Markdown
			MaritalStatus base.Markdown
			Languages     base.Markdown
		}{
			Name:          name,
			Age:           age,
			Gender:        gender,
			Problems:      problems,
			Allergies:     allergies,
			IDDocs:        idDocs,
			Addresses:     addresses,
			MaritalStatus: maritalStatus,
			Languages:     languages,
		}
		t := template.Must(template.New("profile").Parse(patientProfileTemplate))
		_ = t.Execute(w, templateData)
	}
}

// LastVisitHandlerFunc returns a function that renders a patient's last visit
func LastVisitHandlerFunc(ctx context.Context, clinicalService *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		patientID, err := GetPatientID(ctx, id)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData)
		}
		payload, err := clinicalService.GetFHIRPatient(ctx, patientID)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData)
		}
		patient := payload.Resource
		templateData := struct {
			ID           string
			VisitSummary map[string]interface{}
		}{
			ID:           id,
			VisitSummary: patient.RenderVisitSummary(ctx, clinicalService),
		}
		t := template.Must(template.New("visit").Parse(lastVisitTemplate))
		_ = t.Execute(w, templateData)
	}
}

// FullHistoryHandlerFunc returns a function that renders a patient's last visit
func FullHistoryHandlerFunc(ctx context.Context, clinicalService *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		patientID, err := GetPatientID(ctx, id)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData)
		}
		payload, err := clinicalService.GetFHIRPatient(ctx, patientID)
		if err != nil {
			invalidTemplate := template.Must(template.New("invalid").Parse(invalidLinkTemplate))
			invalidTemplateData := struct{}{}
			_ = invalidTemplate.Execute(w, invalidTemplateData)
		}
		patient := payload.Resource
		templateData := struct {
			ID          string
			FullHistory []map[string]interface{}
		}{
			ID:          id,
			FullHistory: patient.RenderFullHistory(ctx, clinicalService),
		}
		t := template.Must(template.New("history").Parse(fullHistoryTemplate))
		_ = t.Execute(w, templateData)
	}
}

func trimString(inp string, maxLength int) string {
	if len(inp) <= maxLength {
		return inp
	}
	return inp[:maxLength-3] + "..."
}
