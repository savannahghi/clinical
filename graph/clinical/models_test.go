package clinical

import (
	"os"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

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

func TestDummy_SetID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid case",
			args: args{
				id: ksuid.New().String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dummy{}
			d.SetID(tt.args.id)
			if d.ID != tt.args.id {
				t.Errorf("set ID not successful")
				return
			}
		})
	}
}
