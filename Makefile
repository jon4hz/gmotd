.PHONY: standalone pam

standalone:
	go build -o gmotd ./cmd/main.go

pam:
	go build -buildmode=c-shared -o pam_gmotd.so .