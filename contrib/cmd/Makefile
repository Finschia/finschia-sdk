SUB = $(shell ls -d -- */)

all build install clean:
	@for d in $(SUB); do pushd $$d; make $@; popd; done
