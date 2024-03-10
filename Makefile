all: config

config: pkg/domain/config/init.pkl.go

pkg/domain/config/init.pkl.go: pkl/*.pkl
	rm -rf pkg/domain/config/*.pkl.go
	pkl-gen-go pkl/config.pkl
