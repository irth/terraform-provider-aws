package opensearchserverless_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccOpenSearchServerlessSecurityPolicyDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_opensearchserverless_security_policy.test"
	dataSourceName := "data.aws_opensearchserverless_security_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.OpenSearchServerlessEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.OpenSearchServerlessEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSecurityPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPolicyDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "policy", resourceName, "policy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "policy_version", resourceName, "policy_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_modified_date"),
				),
			},
		},
	})
}

func testAccSecurityPolicyDataSourceConfig_basic(rName string) string {
	collection := fmt.Sprintf("collection/%s", rName)
	return fmt.Sprintf(`
resource "aws_opensearchserverless_security_policy" "test" {
  name        = %[1]q
  type        = "encryption"
  description = %[1]q
  policy = jsonencode({
    "Rules" = [
      {
        "Resource" = [
          %[2]q
        ],
        "ResourceType" = "collection"
      }
    ],
    "AWSOwnedKey" = true
  })
}

data "aws_opensearchserverless_security_policy" "test" {
  name = aws_opensearchserverless_security_policy.test.name
  type = "encryption"
}
`, rName, collection)
}
