GTASKW_VERSION=latest
# GTASKW_VERSION=v0.0.0-unpublished

tasks: main.go
	# needs to be in go.mod to work against a local go-taskwarrior
	# require github.com/msoulier/go-taskwarrior v0.0.0-unpublished
	# replace github.com/msoulier/go-taskwarrior v0.0.0-unpublished => ../go-taskwarrior
	#go clean -modcache
	#cd ../go-taskwarrior && go build
	#go mod edit -replace=github.com/msoulier/go-taskwarrior@v0.0.0-unpublished=../go-taskwarrior
	go get -d github.com/msoulier/go-taskwarrior@$(GTASKW_VERSION)
	go build

deps: tasks
	./tasks -depgraph > deps.dot
	dot -Tsvg deps.dot -o deps.svg
	inkscape deps.svg

clean:
	go clean -modcache
	rm -f tasks deps.dot deps.svg
