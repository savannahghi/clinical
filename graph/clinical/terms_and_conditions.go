package clinical

import (
	"html/template"
	"net/http"

	"github.com/russross/blackfriday/v2"
)

// ConsumerTermsAndConditionsFunc displays the consumer terms and conditions
func ConsumerTermsAndConditionsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("consumerterms").Parse(consumerTermsTemplate))
		md := []byte(consumerTerms)
		output := blackfriday.Run(md)
		consumerTermsOutput := template.HTML(output)
		_ = t.Execute(w, consumerTermsOutput)
	}
}

// ProviderTermsAndConditionsFunc displays the provider terms and conditions
func ProviderTermsAndConditionsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("providerterms").Parse(providerTermsTemplate))
		md := []byte(providerTerms)
		output := blackfriday.Run(md)
		providerTermsOuput := template.HTML(output)
		_ = t.Execute(w, providerTermsOuput)
	}
}
