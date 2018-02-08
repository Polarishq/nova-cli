package source

import "time"

const AppVersion = "0.4.0"

const DefaultNovaURL = "https://api.splunknova.com"

const validateCredsURLPath = "/v1/account"

const eventsIngestPath = "/v1/events"
const metricsIngestPath = "/v1/metrics?type=custom"

const statsSearchPath = "/v1/search/stats"
const eventsSearchPath = "/v1/search/events"

const metricsListPath = "/v1/metrics"

const configFileRelPath = "/.nova"
const httpTimeout = 10 * time.Second

const ingestionBufferSizeBytes = 1000000 // server side max is 1,048,576
const ingestionBufferTimeout = 1 * time.Second

const novaCLISourcePrefix = "nova-cli-"
