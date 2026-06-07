package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang-jwt/jwt/v5"
)

func TestJwtGenerator_GenerateToken(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		keyPath      string
		inputContent jwt.MapClaims

		expectedOutput string
		expectedErr    error
	}{
		{
			name:    "valid key path",
			keyPath: filepath.FromSlash("./private.test.pem"),
			inputContent: jwt.MapClaims{
				"id":   "1234",
				"name": "John",
			},
			expectedOutput: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQiLCJuYW1lIjoiSm9obiJ9.jWvCQ3JQ0sKPQaQx55iaDA9RtJfa3lFBs7fQp1W6nB_mUrBgQ67naDfnSVFEaQFvY8TyhAy-ivjoqhgirAmfJqOEgVvt7Hm2Rsm4wUHG2F8bJjUX5rgmv8gOEfQT4KkxBZ9yhjUjnKKzluSyvHhHMTzZRaTpiC8tIT3Vf717wxVdIPFnksNm6wkVQzmc_P444im8cElDgKJqJdcdSCu2A8DWpFZkk1kgGXOU28dnWvGqCyASKftaALqAfsUAoFVstFEb5dtZ0VZQrQwaIMPgZAs_X4EadHj0fev8RpLVQBjlsAcJct5rwW7voYr6WSYa5kHUmWZvqw3f0iFArfLsjg",
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testGen, err := NewJWTGenerator(tc.keyPath)
			assert.Equal(t, err, tc.expectedErr)
			res, err := testGen.GenerateToken(tc.inputContent)
			assert.Equal(t, err, tc.expectedErr)
			assert.Equal(t, res, tc.expectedOutput)
		})
	}
}
