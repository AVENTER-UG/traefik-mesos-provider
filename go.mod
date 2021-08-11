module github.com/AVENTER-UG/traefik-mesos-provider

go 1.16

require (
	github.com/m7shapan/njson v1.0.4
	github.com/traefik/genconf v0.1.0
	github.com/traefik/paerser v0.1.4
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
)
