IMAGE=openapitools/openapi-generator-cli:v4.3.1

validate:
	@docker run --rm \
		-v $$PWD:/specs \
		--entrypoint bash \
		${IMAGE} \
		-c 'cp -r /specs/specs/7.0 /tmp/ && \
			(cd /specs && ./generate-spec 7.0 /tmp/7.0/spec.yml) && \
			docker-entrypoint.sh validate -i /tmp/7.0/spec.yml'

GENERATOR=go
go: generate

python: GENERATOR=python
python: generate

ruby: GENERATOR=ruby
ruby: generate

generate:
	@mkdir -p output
	@docker run --rm \
		-v $$PWD:/specs \
		-v $$PWD/output:/output \
		--entrypoint bash \
		${IMAGE} \
		-c 'cp -r /specs/specs/7.0 /tmp/ && \
			(cd /specs && ./generate-spec 7.0 /tmp/7.0/spec.yml) && \
			docker-entrypoint.sh generate -i /tmp/7.0/spec.yml \
				--global-property models,modelTests=false \
				--generator-name ${GENERATOR} --output /output/${GENERATOR}'
