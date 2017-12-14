package src

import "time"

const AppVersion = "0.3.0"

const defaultNovaURL = "https://api.splunknova.com"

const validateCredsURLPath = "/v1/account"
const eventsURLPath = "/v1/events"

const configFileRelPath = "/.nova"
const httpTimeout = 10 * time.Second

const ingestionBufferSizeBytes = 1000000 // server side max is 1,048,576
const ingestionBufferTimeout = 1 * time.Second

const defaultSearchResultsCount = "20"
const novaCLISourcePrefix = "nova-cli-"
