SHELL := bash
.SHELLFLAGS = -ec

export BOOTSTRAPPER_IMG ?= quay.io/amisevsk/workspace-bootstrapper:dev
export SYNC_IMG ?= quay.io/amisevsk/workspace-sync:dev

### help: print this message
help: Makefile
	@{ \
		echo 'Available rules:' ;\
		sed -n 's/^### /    /p' $< | awk 'BEGIN { FS=":" } { printf "%-22s -%s\n", $$1, $$2 }' ;\
		echo '' ;\
		echo 'Supported env vars:' ;\
		echo '    BOOTSTRAPPER_IMG   - image tag for workspace-bootstrapper' ;\
		echo '    SYNC_IMG           - image tag for workspace-sync' ;\
	}

### build: build and push bootstrapper and sync containers
build: build_bootstrapper build_sync

### build_bootstrapper: build and push workspace bootstrapper image
build_bootstrapper:
	docker build -t $(BOOTSTRAPPER_IMG) .
	docker push $(BOOTSTRAPPER_IMG)

### build_sync: build and push workspace sync image
build_sync:
	docker build -t $(SYNC_IMG) -f ./devworkspace-sync/build/Dockerfile .
	docker push $(SYNC_IMG)
