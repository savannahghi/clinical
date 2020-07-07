package clinical

import (
	"html/template"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// ConsumerTermsAndConditionsFunc displays the consumer terms and conditions
func ConsumerTermsAndConditionsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("consumerterms").Parse(consumerTermsTemplate))
		md := []byte(consumerTerms)
		output := blackfriday.Run(md)
		htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(output)
		consumerTermsOutput := template.HTML(htmlSafe)
		_ = t.Execute(w, consumerTermsOutput)
	}
}

// ProviderTermsAndConditionsFunc displays the provider terms and conditions
func ProviderTermsAndConditionsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("providerterms").Parse(providerTermsTemplate))
		md := []byte(providerTerms)
		output := blackfriday.Run(md)
		htmlSafe := bluemonday.UGCPolicy().SanitizeBytes(output)
		providerTermsOuput := template.HTML(htmlSafe)
		_ = t.Execute(w, providerTermsOuput)
	}
}
