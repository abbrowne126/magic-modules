package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/pubsub"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccPubsubSubscription_emptyTTL(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_emptyTTL(topic, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_emptyRetryPolicy(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_emptyRetryPolicy(topic, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_basic(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_basic(topic, subscription, "bar", 20, false),
			},
			{
				ResourceName:            "google_pubsub_subscription.foo",
				ImportStateId:           subscription,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubSubscription_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_basic(topic, subscriptionShort, "bar", 20, false),
			},
			{
				ResourceName:            "google_pubsub_subscription.foo",
				ImportStateId:           subscriptionShort,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccPubsubSubscription_basic(topic, subscriptionShort, "baz", 30, true),
			},
			{
				ResourceName:            "google_pubsub_subscription.foo",
				ImportStateId:           subscriptionShort,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubSubscription_push(t *testing.T) {
	t.Parallel()

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-foo-%s", acctest.RandString(t, 10))
	saAccount := fmt.Sprintf("tf-test-pubsub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_push(topicFoo, saAccount, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_pushNoWrapper(t *testing.T) {
	t.Parallel()

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-foo-%s", acctest.RandString(t, 10))
	saAccount := fmt.Sprintf("tf-test-pubsub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_pushNoWrapper(topicFoo, saAccount, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_pushNoWrapperEmpty(t *testing.T) {
	t.Parallel()

	topicFoo := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-foo-%s", acctest.RandString(t, 10))
	saAccount := fmt.Sprintf("tf-test-pubsub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_pushNoWrapperEmpty(topicFoo, saAccount, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionBigQuery_update(t *testing.T) {
	t.Parallel()

	dataset := fmt.Sprintf("tftestdataset%s", acctest.RandString(t, 10))
	table := fmt.Sprintf("tf-test-table-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscriptionShort, false, ""),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscriptionShort, true, ""),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionBigQuery_serviceAccount(t *testing.T) {
	t.Parallel()

	dataset := fmt.Sprintf("tftestdataset%s", acctest.RandString(t, 10))
	table := fmt.Sprintf("tf-test-table-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com",
			Role:   "roles/bigquery.dataEditor",
		},
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com",
			Role:   "roles/bigquery.metadataViewer",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscriptionShort, false, "bq-test-sa"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscriptionShort, true, ""),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscriptionShort, true, "bq-test-sa2"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionCloudStorage_updateText(t *testing.T) {
	t.Parallel()

	bucket := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "", "", "", 0, "", 0, "", "text"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "pre-", "-suffix", "YYYY-MM-DD/hh_mm_ssZ", 1000, "300s", 1000, "", "text"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionCloudStorage_updateAvro(t *testing.T) {
	t.Parallel()

	bucket := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "", "", "", 0, "", 0, "", "avro"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "pre-", "-suffix", "YYYY-MM-DD/hh_mm_ssZ", 1000, "300s", 1000, "", "avro"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionCloudStorage_emptyAvroConfig(t *testing.T) {
	t.Parallel()

	bucket := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "pre-", "-suffix", "YYYY-MM-DD/hh_mm_ssZ", 1000, "300s", 1000, "", "empty-avro"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscriptionCloudStorage_serviceAccount(t *testing.T) {
	t.Parallel()

	bucket := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "", "", "", 0, "", 0, "gcs-test-sa", "text"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "pre-", "-suffix", "YYYY-MM-DD/hh_mm_ssZ", 1000, "300s", 1000, "", "text"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscriptionShort, "", "", "", 0, "", 0, "gcs-test-sa2", "avro"),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Context: hashicorp/terraform-provider-google#4993
// This test makes a call to GET an subscription before it is actually created.
// The PubSub API negative-caches responses so this tests we are
// correctly polling for existence post-creation.
func TestAccPubsubSubscription_pollOnCreate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-topic-foo-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create only the topic
				Config: testAccPubsubSubscription_topicOnly(topic),
				// Read from non-existent subscription created in next step
				// so API negative-caches result
				Check: testAccCheckPubsubSubscriptionCache404(t, subscription),
			},
			{
				// Create the subscription - if the polling fails,
				// the test step will fail because the read post-create
				// will have removed the resource from state.
				Config: testAccPubsubSubscription_pollOnCreate(topic, subscription),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscription,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitPubsubSubscription_IgnoreMissingKeyInMap(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		Key                string
		ExpectDiffSuppress bool
	}{
		"missing key in map": {
			Old:                "",
			New:                "v1",
			Key:                "x-goog-version",
			ExpectDiffSuppress: true,
		},
		"different values": {
			Old:                "v1",
			New:                "v2",
			Key:                "x-goog-version",
			ExpectDiffSuppress: false,
		},
		"same values": {
			Old:                "v1",
			New:                "v1",
			Key:                "x-goog-version",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if pubsub.IgnoreMissingKeyInMap(tc.Key)("push_config.0.attributes."+tc.Key, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestAccPubsubSubscription_filter(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscription_filter(topic, subscriptionShort, "attributes.foo = \\\"bar\\\""),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "filter", "attributes.foo = \"bar\""),
				),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSubscription_filter(topic, subscriptionShort, ""),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "filter", ""),
				),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubSubscription_javascriptUdfUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscriptionShort := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))
	functionName1 := "filter_falsy"
	functionName2 := "passthrough"
	code1 := "function filter_falsy(message, metadata) {\n  return message ? message : null\n}\n"
	code2 := "function passthrough(message, metadata) {\n    return message\n}\n"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			// Initial transform
			{
				Config: testAccPubsubSubscription_javascriptUdfSettings(topic, subscriptionShort, functionName1, code1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.function_name", functionName1),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.code", code1),
				),
			},
			// Bare transform
			{
				Config: testAccPubsubSubscription_javascriptUdfSettings(topic, subscriptionShort, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.function_name", ""),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.code", ""),
				),
			},
			// Destroy transform
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Two transforms
			{
				Config: testAccPubsubSubscription_javascriptUdfSettings(topic, subscriptionShort, functionName1, code1) + "\n" + testAccPubsubSubscription_javascriptUdfSettings(topic, subscriptionShort, functionName2, code2),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.function_name", functionName1),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.code", code1),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.1.function_name", functionName2),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.1.code", code2),
				),
			},
			{
				// Remove non-required field
				Config: testAccPubsubSubscription_javascriptUdfSettings_noEnabled(topic, subscriptionShort, functionName1, code1),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.function_name", functionName1),
					resource.TestCheckResourceAttr("google_pubsub_subscription.foo", "message_transforms.0.code", code1),
				),
			},
			{
				ResourceName:      "google_pubsub_subscription.foo",
				ImportStateId:     subscriptionShort,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubSubscription_javascriptUdfSettings(topic, subscription, functionName, code string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id
	message_transforms {
		{
			javascript_udf {
				function_name = %s,
				code = %s
			}
			disabled = false
		}
  }
}
`, topic, subscription, functionName, code)
}

func testAccPubsubSubscription_javascriptUdfSettings_noEnabled(topic, subscription, functionName, code string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id
	message_transforms {
		{
			javascript_udf {
				function_name = %s,
				code = %s
			}
		}
  }
}
`, topic, subscription, functionName, code)
}

func testAccPubsubSubscription_emptyTTL(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id

  message_retention_duration = "1200s"
  retain_acked_messages      = true
  ack_deadline_seconds       = 20
  expiration_policy {
    ttl = ""
  }
  enable_message_ordering    = false
}
`, topic, subscription)
}

func testAccPubsubSubscription_emptyRetryPolicy(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id

  retry_policy {
  }
}
`, topic, subscription)
}

func testAccPubsubSubscription_push(topicFoo, saAccount, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_service_account" "pub_sub_service_account" {
  account_id = "%s"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/projects.topics.publish"

    members = [
      "serviceAccount:${google_service_account.pub_sub_service_account.email}",
    ]
  }
}

resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name                 = "%s"
  topic                = google_pubsub_topic.foo.name
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = google_service_account.pub_sub_service_account.email
    }
  }
}
`, saAccount, topicFoo, subscription)
}

func testAccPubsubSubscription_pushNoWrapper(topicFoo, saAccount, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_service_account" "pub_sub_service_account" {
  account_id = "%s"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/projects.topics.publish"

    members = [
      "serviceAccount:${google_service_account.pub_sub_service_account.email}",
    ]
  }
}

resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name                 = "%s"
  topic                = google_pubsub_topic.foo.name
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = google_service_account.pub_sub_service_account.email
    }
    no_wrapper {
      write_metadata = true
    }
  }
}
`, saAccount, topicFoo, subscription)
}

func testAccPubsubSubscription_pushNoWrapperEmpty(topicFoo, saAccount, subscription string) string {
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_service_account" "pub_sub_service_account" {
  account_id = "%s"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/projects.topics.publish"

    members = [
      "serviceAccount:${google_service_account.pub_sub_service_account.email}",
    ]
  }
}

resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name                 = "%s"
  topic                = google_pubsub_topic.foo.name
  ack_deadline_seconds = 10
  push_config {
    push_endpoint = "https://${data.google_project.project.project_id}.appspot.com"
    oidc_token {
      service_account_email = google_service_account.pub_sub_service_account.email
    }
    no_wrapper {
      write_metadata = false
    }
  }
}
`, saAccount, topicFoo, subscription)
}

func testAccPubsubSubscription_basic(topic, subscription, label string, deadline int, exactlyOnceDelivery bool) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name   = "%s"
  topic  = google_pubsub_topic.foo.id
  filter = "attributes.foo = \"bar\""
  labels = {
    foo = "%s"
  }
  retry_policy {
    minimum_backoff = "60.0s"
  }
  ack_deadline_seconds = %d
  enable_exactly_once_delivery = %t
}
`, topic, subscription, label, deadline, exactlyOnceDelivery)
}

func testAccPubsubSubscriptionBigQuery_basic(dataset, table, topic, subscription string, useTableSchema bool, serviceAccountId string) string {
	serviceAccountEmailField := ""
	serviceAccountResource := ""
	tfDependencies := ""
	if serviceAccountId != "" {
		serviceAccountResource = fmt.Sprintf(`
resource "google_service_account" "bq_write_service_account" {
  account_id   = "%s"
  display_name = "BQ Write Service Account"
}

resource "google_project_iam_member" "bigquery_metadata_viewer" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.metadataViewer"
  member  = "serviceAccount:${google_service_account.bq_write_service_account.email}"
}

resource "google_project_iam_member" "bigquery_data_editor" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:${google_service_account.bq_write_service_account.email}"
}`, serviceAccountId)
		serviceAccountEmailField = "service_account_email = google_service_account.bq_write_service_account.email"
		tfDependencies = `    google_project_iam_member.bigquery_metadata_viewer,
    google_project_iam_member.bigquery_data_editor,
    time_sleep.wait_30_seconds,`
	} else {
		tfDependencies = "    time_sleep.wait_30_seconds,"
	}
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "time_sleep" "wait_30_seconds" {
  create_duration = "30s"
}

%s

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema = <<EOF
[
  {
    "name": "data",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The data"
  }
]
EOF
}

resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name   = "%s"
  topic  = google_pubsub_topic.foo.id

  bigquery_config {
    table = "${google_bigquery_table.test.project}.${google_bigquery_table.test.dataset_id}.${google_bigquery_table.test.table_id}"
    use_table_schema = %t
    %s
  }

  depends_on = [
    %s
  ]
}
	`, serviceAccountResource, dataset, table, topic, subscription, useTableSchema, serviceAccountEmailField, tfDependencies)
}

func testAccPubsubSubscriptionCloudStorage_basic(bucket, topic, subscription, filenamePrefix, filenameSuffix, filenameDatetimeFormat string, maxBytes int, maxDuration string, maxMessages int, serviceAccountId, outputFormat string) string {
	filenamePrefixString := ""
	if filenamePrefix != "" {
		filenamePrefixString = fmt.Sprintf(`filename_prefix = "%s"`, filenamePrefix)
	}
	filenameSuffixString := ""
	if filenameSuffix != "" {
		filenameSuffixString = fmt.Sprintf(`filename_suffix = "%s"`, filenameSuffix)
	}
	filenameDatetimeString := ""
	if filenameDatetimeFormat != "" {
		filenameDatetimeString = fmt.Sprintf(`filename_datetime_format = "%s"`, filenameDatetimeFormat)
	}
	maxBytesString := ""
	if maxBytes != 0 {
		maxBytesString = fmt.Sprintf(`max_bytes = %d`, maxBytes)
	}
	maxDurationString := ""
	if maxDuration != "" {
		maxDurationString = fmt.Sprintf(`max_duration = "%s"`, maxDuration)
	}
	maxMessagesString := ""
	if maxMessages != 0 {
		maxMessagesString = fmt.Sprintf(`max_messages = %d`, maxMessages)
	}

	serviceAccountEmailField := ""
	serviceAccountResource := ""
	if serviceAccountId != "" {
		serviceAccountResource = fmt.Sprintf(`
resource "google_service_account" "storage_write_service_account" {
  account_id   = "%s"
  display_name = "Write Service Account"
}

resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.test.name
  role   = "roles/storage.admin"
	member = "serviceAccount:${google_service_account.storage_write_service_account.email}"
}

resource "google_project_iam_member" "editor" {
	project = data.google_project.project.project_id
	role   = "roles/bigquery.dataEditor"
	member = "serviceAccount:${google_service_account.storage_write_service_account.email}"
}`, serviceAccountId)
		serviceAccountEmailField = "service_account_email = google_service_account.storage_write_service_account.email"
	} else {
		serviceAccountResource = fmt.Sprintf(`
resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.test.name
  role   = "roles/storage.admin"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}`)
	}
	outputFormatString := ""
	if outputFormat == "avro" {
		outputFormatString = `
  avro_config {
    write_metadata = true
    use_topic_schema = true
  }
`
	} else if outputFormat == "empty-avro" {
		outputFormatString = `avro_config {}`
	}
	return fmt.Sprintf(`
data "google_project" "project" { }

resource "google_storage_bucket" "test" {
  name = "%s"
  location = "US"
}

%s

resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name   = "%s"
  topic  = google_pubsub_topic.foo.id

  cloud_storage_config {
    bucket = "${google_storage_bucket.test.name}"
    %s
    %s
    %s
    %s
    %s
    %s
    %s
    %s
  }

  depends_on = [
    google_storage_bucket.test,
    google_storage_bucket_iam_member.admin,
  ]
}
`, bucket, serviceAccountResource, topic, subscription, filenamePrefixString, filenameSuffixString, filenameDatetimeString, maxBytesString, maxDurationString, maxMessagesString, serviceAccountEmailField, outputFormatString)
}

func testAccPubsubSubscription_topicOnly(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}
`, topic)
}

func testAccPubsubSubscription_pollOnCreate(topic, subscription string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id
}
`, topic, subscription)
}

func TestGetComputedTopicName(t *testing.T) {
	type testData struct {
		project  string
		topic    string
		expected string
	}

	var testCases = []testData{
		{
			project:  "my-project",
			topic:    "my-topic",
			expected: "projects/my-project/topics/my-topic",
		},
		{
			project:  "my-project",
			topic:    "projects/another-project/topics/my-topic",
			expected: "projects/another-project/topics/my-topic",
		},
	}

	for _, testCase := range testCases {
		computedTopicName := pubsub.GetComputedTopicName(testCase.project, testCase.topic)
		if computedTopicName != testCase.expected {
			t.Fatalf("bad computed topic name: %s' => expected %s", computedTopicName, testCase.expected)
		}
	}
}

func testAccCheckPubsubSubscriptionCache404(t *testing.T, subName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		url := fmt.Sprintf("%sprojects/%s/subscriptions/%s", config.PubsubBasePath, envvar.GetTestProjectFromEnv(), subName)
		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err == nil {
			return fmt.Errorf("Expected Pubsub Subscription %q not to exist, was found", resp["name"])
		}
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Got non-404 error while trying to read Pubsub Subscription %q: %v", subName, err)
		}
		return nil
	}
}

func testAccPubsubSubscription_filter(topic, subscription, filter string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
}

resource "google_pubsub_subscription" "foo" {
  name   = "%s"
  topic  = google_pubsub_topic.foo.id
  filter = "%s"
}
`, topic, subscription, filter)
}
