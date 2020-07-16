IMAGE=openapitools/openapi-generator-cli:v4.3.1

validate:
	@docker run --rm \
		-v $$PWD:/specs \
		--entrypoint bash \
		${IMAGE} \
		-c 'cp -r /specs/specs/7.0 /tmp/ && \
			(cd /specs && ./generate-spec 7.0 /tmp/7.0/spec.yml) && \
			docker-entrypoint.sh validate -i /tmp/7.0/spec.yml'
