package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_update(topic, "foo", "bar"),
			},
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topic,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccPubsubTopic_updateWithRegion(topic, "wibble", "wobble", "us-central1"),
			},
			{
				ResourceName:            "google_pubsub_topic.foo",
				ImportStateId:           topic,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubTopic_cmek(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKey(t)
	topicName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-pubsub.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_cmek(topicName, kms.CryptoKey.Name),
			},
			{
				ResourceName:      "google_pubsub_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_schema(t *testing.T) {
	t.Parallel()

	schema1 := fmt.Sprintf("tf-test-schema-%s", acctest.RandString(t, 10))
	schema2 := fmt.Sprintf("tf-test-schema-%s", acctest.RandString(t, 10))
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithSchema(topic, schema1),
			},
			{
				ResourceName:      "google_pubsub_topic.bar",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithNewSchema(topic, schema2),
			},
			{
				Config: testAccPubsubTopic_updateWithNewSchema(topic, ""),
			},
			{
				ResourceName:      "google_pubsub_topic.bar",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_migration(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.84.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccPubsubTopic_update(topic, "foo", "bar"),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccPubsubTopic_update(topic, "foo", "bar"),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
			{
				ResourceName:             "google_pubsub_topic.foo",
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ImportStateId:            topic,
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateVerifyIgnore:  []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccPubsubTopic_kinesisIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithKinesisIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedKinesisIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_cloudStorageIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithCloudStorageIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedCloudStorageIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_update(topic, key, value string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }
}
`, topic, key, value)
}

func testAccPubsubTopic_updateWithRegion(topic, key, value, region string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    %s = "%s"
  }

  message_storage_policy {
    allowed_persistence_regions = [
      "%s",
    ]
    enforce_in_transit = false
  }
}
`, topic, key, value, region)
}

func testAccPubsubTopic_cmek(topicName, kmsKey string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  name         = "%s"
  kms_key_name = "%s"
}
`, topicName, kmsKey)
}

func testAccPubsubTopic_updateWithSchema(topic, schema string) string {
	return fmt.Sprintf(`
resource "google_pubsub_schema" "foo" {
	name = "%s"
	type = "PROTOCOL_BUFFER"
  definition = "syntax = \"proto3\";\nmessage Results {\nstring f1 = 1;\n}"
}

resource "google_pubsub_topic" "bar" {
  name = "%s"
	schema_settings {
    schema = google_pubsub_schema.foo.id
    encoding = "BINARY"
  }
}
`, schema, topic)
}

func testAccPubsubTopic_updateWithNewSchema(topic, schema string) string {
	if schema != "" {
		return fmt.Sprintf(`
resource "google_pubsub_schema" "foo" {
	name = "%s"
	type = "PROTOCOL_BUFFER"
	definition = "syntax = \"proto3\";\nmessage Results {\nstring f1 = 1;\n}"
}

resource "google_pubsub_topic" "bar" {
  name = "%s"
	schema_settings {
    schema = google_pubsub_schema.foo.id
    encoding = "JSON"
  }
}
`, schema, topic)
	} else {
		return fmt.Sprintf(`
		resource "google_pubsub_topic" "bar" {
			name = "%s"
		}
		`, topic)
	}
}

func testAccPubsubTopic_updateWithKinesisIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual AWS resources for the test to pass.
  ingestion_data_source_settings {
    aws_kinesis {
        stream_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name"
        consumer_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name/consumer/consumer-1:1111111111"
        aws_role_arn = "arn:aws:iam::111111111111:role/fake-role-name"
        gcp_service_account = "fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedKinesisIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual AWS resources for the test to pass.
  ingestion_data_source_settings {
    aws_kinesis {
        stream_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/updated-fake-stream-name"
        consumer_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/updated-fake-stream-name/consumer/consumer-1:1111111111"
        aws_role_arn = "arn:aws:iam::111111111111:role/updated-fake-role-name"
        gcp_service_account = "updated-fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithCloudStorageIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
    cloud_storage {
        bucket = "test-bucket"
        text_format {
            delimiter = " "
        }
        minimum_object_create_time = "2024-01-01T00:00:00Z"
        match_glob = "foo/**"
    }
    platform_logs_settings {
        severity = "WARNING"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedCloudStorageIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
    cloud_storage {
        bucket = "updated-test-bucket"
        avro_format {}
        minimum_object_create_time = "2024-02-02T00:00:00Z"
        match_glob = "bar/**"
    }
    platform_logs_settings {
        severity = "ERROR"
    }
  }
}
`, topic)
}

func TestAccPubsubTopic_azureEventHubsIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithAzureEventHubsIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedAzureEventHubsIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_updateWithAzureEventHubsIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
  	azure_event_hubs {
		resource_group = "azure-ingestion-resource-group"
		namespace = "azure-ingestion-namespace"
		event_hub = "azure-ingestion-event-hub"
		client_id = "aZZZZZZZ-YYYY-HHHH-GGGG-abcdef569123"
		tenant_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456789123"
		subscription_id = "bXXXXXXX-YYYY-HHHH-GGGG-123456789123"
		gcp_service_account = "fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedAzureEventHubsIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
  	azure_event_hubs {
		resource_group = "ingestion-resource-group"
		namespace = "ingestion-namespace"
		event_hub = "ingestion-event-hub"
		client_id = "aZZZZZZZ-YYYY-HHHH-GGGG-abcdef123456"
		tenant_id = "0XXXXXXX-YYYY-HHHH-GGGG-123456123456"
		subscription_id = "bXXXXXXX-YYYY-HHHH-GGGG-123456123456"
		gcp_service_account = "fake-account@new-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func TestAccPubsubTopic_awsMskIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithAwsMskIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedAwsMskIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_updateWithAwsMskIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
  	aws_msk {
		cluster_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name"
		topic = "test-topic"
		aws_role_arn = "arn:aws:iam::111111111111:role/fake-role-name"
		gcp_service_account = "fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedAwsMskIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
  	aws_msk {
		cluster_arn = "arn:aws:kinesis:us-west-2:111111111111:stream/fake-stream-name"
		topic = "test-topic"
		aws_role_arn = "arn:aws:iam::111111111111:role/fake-role-name"
		gcp_service_account = "updated-fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func TestAccPubsubTopic_confluentCloudIngestionUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_updateWithConfluentCloudIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_updateWithUpdatedConfluentCloudIngestionSettings(topic),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubTopic_updateWithConfluentCloudIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
	confluent_cloud {
		bootstrap_server = "test.us-west2.gcp.confluent.cloud:1111"
		cluster_id = "1234"
		topic = "test-topic"
		identity_pool_id = "test-identity-pool-id"
		gcp_service_account = "fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}

func testAccPubsubTopic_updateWithUpdatedConfluentCloudIngestionSettings(topic string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"

  # Outside of automated terraform-provider-google CI tests, these values must be of actual Cloud Storage resources for the test to pass.
  ingestion_data_source_settings {
	confluent_cloud {
		bootstrap_server = "test.us-west2.gcp.confluent.cloud:1111"
		cluster_id = "1234"
		topic = "test-topic"
		identity_pool_id = "test-identity-pool-id"
		gcp_service_account = "updated-fake-service-account@fake-gcp-project.iam.gserviceaccount.com"
    }
  }
}
`, topic)
}
func TestAccPubsubTopic_javascriptUdfUpdate(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	functionName1 := "filter_falsy"
	functionName2 := "passthrough"
	code1 := "function filter_falsy(message, metadata) {\n  return message ? message : null;\n}\n"
	code2 := "function passthrough(message, metadata) {\n    return message;\n}\n"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			// Initial transform
			{
				Config: testAccPubsubTopic_javascriptUdfSettings(topic, functionName1, code1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.function_name", functionName1),
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.code", code1),
				),
			},
			// Bare transform
			{
				Config: testAccPubsubTopic_javascriptUdfSettings(topic, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.function_name", ""),
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.code", ""),
				),
			},
			// Destroy transform
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Two transforms
			{
				Config: testAccPubsubTopic_javascriptUdfSettings(topic, functionName1, code1) + "\n" + testAccPubsubTopic_javascriptUdfSettings(topic, functionName2, code2),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.function_name", functionName1),
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.0.code", code1),
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.1.function_name", functionName2),
					resource.TestCheckResourceAttr("google_pubsub_topic.foo", "message_transforms.1.code", code2),
				),
			},
		},
	})
}

func testAccPubsubTopic_javascriptUdfSettings(topic, functionName, code string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"

	message_transforms {
		{
			javascript_udf = {
				function_name = %s,
				code = %s
			}
			disabled = false
		}
	}
}
	`, topic, functionName, code)
}
