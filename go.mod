module github.com/MenD32/Tempest

go 1.23.4

require (
	github.com/MenD32/Shakespeare v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.1
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.130.1
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/MenD32/Shakespeare => ../Shakespeare
