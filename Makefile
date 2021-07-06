SPEC_VERSION ?= 7.0

validate:
	@swagger-cli validate \
		--no-schema \
		specs/${SPEC_VERSION}/spec.yml

bundle: validate
	@swagger-cli bundle \
		--dereference \
		--outfile _gen/${SPEC_VERSION}/spec.json \
		specs/${SPEC_VERSION}/spec.yml

drone:
	drone lint
	drone --server https://drone.grafana.net sign --save grafana/dashboard-spec

.PHONY: validate bundle drone
