build:
	docker build -t generate-resumes .

init:
	mkdir -p outputs inputs
	cp examples/*.toml inputs

run:
	docker run -v "$(shell pwd)/outputs:/outputs" -v "$(shell pwd)/inputs:/inputs" generate-resumes /inputs/$(filename).toml -o /outputs

clean:
	rm -rf outputs
	rm -rf inputs

