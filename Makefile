SPEC_VERSION ?= 7.0

validate:
	@swagger-cli validate \
		--no-schema \
		specs/${SPEC_VERSION}/spec.yml

bundle: validate
	@swagger-cli bundle \
		--dereference \
		--outfile bundle/${SPEC_VERSION}/spec.json \
		specs/${SPEC_VERSION}/spec.yml

GENERATOR_IMAGE = openapitools/openapi-generator-cli:v4.3.1
GENERATOR ?= go

generate: bundle
	@docker run --rm \
		-v $$PWD/bundle:/bundle \
		-v $$PWD/_gen:/gen \
		${GENERATOR_IMAGE} \
		generate -i /bundle/${SPEC_VERSION}/spec.json \
			--global-property models,modelTests=false \
			--generator-name ${GENERATOR} --output /gen/${SPEC_VERSION}/${GENERATOR}

.PHONY: validate bundle generate
