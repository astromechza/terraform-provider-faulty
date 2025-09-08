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

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "faulty" {	
  required_boolean = true
}
resource "faulty_example" "test" {
  required_boolean = true
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"faulty_example.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"faulty_example.test",
						tfjsonpath.New("required_boolean"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: `
provider "faulty" {	
  required_boolean = true
}
resource "faulty_example" "test" {
  required_boolean = true
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"faulty_example.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"faulty_example.test",
						tfjsonpath.New("required_boolean"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccExampleResource_fail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "faulty_example" "test" {
  required_boolean = true
}
`,
				ExpectError: regexp.MustCompile("Invalid provider configuration"),
			},
			{
				Config: `
provider "faulty" {
}

resource "faulty_example" "test" {
  required_boolean = true
}
`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config: `
provider "faulty" {
  required_boolean = false
}

resource "faulty_example" "test" {
  required_boolean = true
}
`,
				ExpectError: regexp.MustCompile("required_boolean_not_true"),
			},
			{
				Config: `
provider "faulty" {
  required_boolean = true
}

resource "faulty_example" "test" {
  required_boolean = false
}
`,
				ExpectError: regexp.MustCompile("required_boolean_not_true"),
			},
		},
	})
}
