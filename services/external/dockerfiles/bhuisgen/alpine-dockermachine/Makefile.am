IMAGE = bhuisgen/alpine-dockermachine

clean-local:
	docker rmi $(IMAGE):$(TAG) $(CLEAN_OPTIONS) || true

build-local:
	docker build -t $(IMAGE):$(TAG) $(BUILD_OPTIONS) .

pull-local:
	docker pull $(IMAGE):$(TAG) $(PULL_OPTIONS)
