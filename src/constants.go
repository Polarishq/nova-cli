package src

import "time"

const AppVersion = "0.3.1"

const DefaultNovaURL = "https://api.splunknova.com"

const validateCredsURLPath = "/v1/account"
const eventsURLPath = "/v1/events"

const metricsURLSearchPath = "/v1/metrics"
const metricsURLIngestPath = "/v1/metrics?type=custom"

const configFileRelPath = "/.nova"
const httpTimeout = 20 * time.Second

const ingestionBufferSizeBytes = 1000000 // server side max is 1,048,576
const ingestionBufferTimeout = 1 * time.Second

const defaultSearchResultsCount = "20"
const novaCLISourcePrefix = "nova-cli-"
