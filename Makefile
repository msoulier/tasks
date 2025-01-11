tasks: main.go
	go clean -modcache
	cd ../go-taskwarrior && go build
	go mod edit -replace=github.com/msoulier/go-taskwarrior@v0.0.0-unpublished=../go-taskwarrior
	go get -d github.com/msoulier/go-taskwarrior@v0.0.0-unpublished
	go build

deps: tasks
	./tasks -depgraph > deps.dot
	dot -Tsvg deps.dot -o deps.svg
	inkscape deps.svg

clean:
	go clean -modcache
	rm -f tasks deps.dot deps.svg
