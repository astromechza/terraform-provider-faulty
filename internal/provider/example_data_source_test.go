// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleDataSource_valid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `
provider "faulty" {
  required_boolean = true
}

data "faulty_example" "test" {
  required_boolean = true
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.faulty_example.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccExampleDataSource_failures(t *testing.T) {
	for _, tc := range []struct {
		name        string
		config      string
		expectError *regexp.Regexp
	}{
		{
			name:        "provider not present",
			expectError: regexp.MustCompile("Invalid provider configuration"),
			config: `
data "faulty_example" "test" {
  required_boolean = true
}
`,
		},
		{
			name:        "required boolean false",
			expectError: regexp.MustCompile("required_boolean_not_true"),
			config: `
provider "faulty" {
  required_boolean = false
}

data "faulty_example" "test" {
  required_boolean = true
}
`,
		},
		{
			name:        "required boolean not set",
			expectError: regexp.MustCompile("Missing required argument"),
			config: `
provider "faulty" {
}

data "faulty_example" "test" {
  required_boolean = true
}
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: tc.expectError,
					},
				},
			})
		})
	}
}
