package provider

import (
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func setupHelmRelease(t *testing.T) {
	out, err := exec.Command("helm", "list", "-A", "-o", "json").Output()
	data := make([]map[string]interface{}, 0)
	json.Unmarshal(out, &data)

	t.Log(data)

	release_exists := false
	for _, release := range data {
		if release["name"] == "hello-world" && release["namespace"] == "terraform-provider-torb-testing" {
			release_exists = true
		}
	}

	if !release_exists {
		if err != nil {
			t.Fatalf("Error setting up, checking if release already exists: %s", err)
		}

		_, err = exec.Command("helm", "repo", "add", "cloudecho", "https://cloudecho.github.io/charts/").Output()
		if err != nil {
			t.Fatalf("Error setting up, repo add: %s", err)
		}
		_, err = exec.Command("helm", "repo", "update").Output()
		if err != nil {
			t.Fatalf("Error setting up, repo update: %s", err)
		}
		_, err = exec.Command("helm", "install", "hello-world", "cloudecho/hello", "--namespace", "terraform-provider-torb-testing", "--version=0.1.2", "--create-namespace").Output()
		if err != nil {
			t.Fatalf("Error setting up, install: %s", err)
		}
	}

}

func cleanup(t *testing.T) {
	out, err := exec.Command("helm", "uninstall", "hello-world", "--namespace", "terraform-provider-torb-testing").Output()
	if err != nil {
		t.Fatalf("Error setting up helm release: %s", out)
	}
}

func TestAccHelmReleaseDataSource(t *testing.T) {
	setupHelmRelease(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			t.FailNow()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccHelmReleaseDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.helm_release.test", "values", ""),
				),
			},
		},
	})

	cleanup(t)
}

const testAccHelmReleaseDataSourceConfig = `
data "helm_release" "test" {
	namespace = "terraform-provider-torb-testing"
	release_name = "hello-world"
}
`
