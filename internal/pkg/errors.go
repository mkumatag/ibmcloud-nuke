package pkg

import "fmt"

var (
	ErrorResourceTypeMatch = fmt.Errorf("resource type doesn't match")

	ErrorMissingRegion = fmt.Errorf("region is not specified in resouce or globally")

	ErrorResourceNotFound = fmt.Errorf("resource not found")
)
