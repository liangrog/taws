// Managing AWS login session
package utils

import (
	_ "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

const (
	errMissingRegion     = "Missing environment variable AWS_REGION or option --region"
	errMissingCredential = "Missing environment variable AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY, alternatively you can user -profile option"
)

// Global AWS session
// Also can be used for testing stud
// Default empty session, you need to call
// GetSession on each command to get
// The right setting from persistent flags
var AwsSess *session.Session = session.Must(session.NewSession())

// Get AWS via ENV or profile options
func GetSession(profile, region string) *session.Session {
	// Region is required to be set
	_, rgnOk := os.LookupEnv("AWS_REGION")
	if !rgnOk && len(region) == 0 {
		panic(errMissingRegion)
	}

	options := session.Options{}

	// If region is set via option
	if len(region) > 0 {
		options.Config = aws.Config{Region: aws.String(region)}
	}

	// If profile is set via option
	if len(profile) > 0 {
		options.Profile = profile
	}

	// We require either key and secret pack or profile defined
	_, keyOk := os.LookupEnv("AWS_ACCESS_KEY_ID")
	_, secretOk := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	_, pflOk := os.LookupEnv("AWS_PROFILE")

	if (!keyOk || !secretOk) &&
		(len(profile) == 0 && !pflOk) {
		panic(errMissingCredential)
	}

	// Assign session to global
	AwsSess = session.Must(session.NewSessionWithOptions(options))

	return AwsSess
}
