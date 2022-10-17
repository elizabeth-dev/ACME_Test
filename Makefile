proto:
	@sh -c "'$(CURDIR)/scripts/proto.sh'"

package package.user:
	@sh -c "'$(CURDIR)/scripts/package.user.sh'"

run run.compose:
	@sh -c "'$(CURDIR)/scripts/run.compose.sh'"

.PHONY: mocks
mocks:
	@sh -c "'$(CURDIR)/scripts/mocks.sh'"

test:
	@sh -c "'$(CURDIR)/scripts/test.sh'"
