#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/line/link/app

sim-link-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -v -timeout 24h

sim-link-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.linkd/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.linkd/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-link-fast:
	@echo "Running quick link simulation. This may take several minutes..."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-link-import-export: runsim
	@echo "Running link import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 25 5 TestAppImportExport

sim-link-simulation-after-import: runsim
	@echo "Running link simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 25 5 TestAppSimulationAfterImport

sim-link-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.linkd/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -g ${HOME}/.linkd/config/genesis.json 400 5 TestFullAppSimulation

sim-link-multi-seed: runsim
	@echo "Running multi-seed link simulation. This may take awhile!"
	$(GOPATH)/bin/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 400 5 TestFullAppSimulation

sim-link-multi-seed-short: runsim
	@echo "Running short multi-seed link simulation. This may take awhile!"
	$(GOPATH)/bin/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 40 5 TestFullAppSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-link-benchmark:
	@echo "Running link benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-link-profile:
	@echo "Running link benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-link-nondeterminism sim-link-custom-genesis-fast sim-link-fast sim-link-import-export \
	sim-link-simulation-after-import sim-link-custom-genesis-multi-seed sim-link-multi-seed \
	sim-benchmark-invariants sim-link-benchmark sim-link-profile
