module github.com/bsdlp/delete-twitter-data

go 1.15

require (
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.6.0
	github.com/kelseyhightower/envconfig v1.4.0
)

replace github.com/kelseyhightower/envconfig v1.4.0 => github.com/bsdlp/envconfig v1.5.0
