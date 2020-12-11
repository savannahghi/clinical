package clinical

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpisodeOfCareStatus(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        EpisodeOfCareStatus
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    EpisodeOfCareStatusPlanned,
			convert: EpisodeOfCareStatusPlanned,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    EpisodeOfCareStatusActive,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    EpisodeOfCareStatusFinished,
			convert: EpisodeOfCareStatusFinished,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)
		})
	}
}

func TestIDDocumentType(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        IDDocumentType
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    IDDocumentTypeNationalID,
			convert: IDDocumentTypeNationalID,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    IDDocumentTypePassport,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    IDDocumentTypePassport,
			convert: IDDocumentTypePassport,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestMaritalStatus(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        MaritalStatus
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    MaritalStatusS,
			convert: MaritalStatusS,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    MaritalStatusU,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    MaritalStatusU,
			convert: MaritalStatusU,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestRelationshipType(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        RelationshipType
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    RelationshipTypeF,
			convert: RelationshipTypeF,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    RelationshipTypeO,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    RelationshipTypeO,
			convert: RelationshipTypeO,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestIdentifierUse(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        IdentifierUse
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    IdentifierUseOld,
			convert: IdentifierUseTemp,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    IdentifierUseTemp,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    IdentifierUseTemp,
			convert: IdentifierUseTemp,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestContactPointSystem(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        ContactPointSystem
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    ContactPointSystemFax,
			convert: ContactPointSystemFax,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    ContactPointSystemEmail,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    ContactPointSystemEmail,
			convert: ContactPointSystemEmail,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestContactPointUse(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        ContactPointUse
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    ContactPointUseHome,
			convert: ContactPointUseHome,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    ContactPointUseMobile,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    ContactPointUseMobile,
			convert: ContactPointUseMobile,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestNameUse(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        NameUse
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    NameUseTemp,
			convert: NameUseTemp,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    NameUseNickname,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    NameUseNickname,
			convert: NameUseNickname,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestAddressType(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        AddressType
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    AddressTypePostal,
			convert: AddressTypePostal,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    AddressTypePostal,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    AddressTypePostal,
			convert: AddressTypePostal,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestAddressUse(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        AddressUse
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    AddressUseWork,
			convert: AddressUseWork,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    AddressUseWork,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    AddressUseWork,
			convert: AddressUseWork,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestCountry(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        Country
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testaddres",
			convert: "testaddress",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    CountryBh,
			convert: CountryMe,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    CountryHn,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    CountryHn,
			convert: CountryHn,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}
