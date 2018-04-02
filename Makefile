TG_BOT_TOKEN=
TG_BOT_SECRET=god
TG_BOT_DOMAIN=
TG_BOT_CERT=cert.pem
TG_BOT_KEY=key.pem
CERT_SUBJ="/C=SA/ST=South Africa/O=Men in black/CN=mib.space"

.PHONY: binary certs

production: binary
	TG_BOT_TOKEN=$(TG_BOT_TOKEN) \
	TG_BOT_SECRET=$(TG_BOT_SECRET) \
	TG_BOT_DOMAIN=$(TG_BOT_DOMAIN) \
	TG_BOT_CERT=$(TG_BOT_CERT) \
	TG_BOT_KEY=$(TG_BOT_KEY) \
	./tgflunky

local: binary
	TG_BOT_TOKEN=$(TG_BOT_TOKEN) \
	TG_BOT_SECRET=$(TG_BOT_SECRET) \
	./tgflunky

debug: binary
	TG_BOT_TOKEN=$(TG_BOT_TOKEN) \
	TG_BOT_SECRET=$(TG_BOT_SECRET) \
	./tgflunky -debug

binary:
	go build -o ./tgflunky ./cmd

certs:
	openssl req -newkey rsa:4096 -sha256 -nodes -keyout key.pem -x509 -days 365 -out cert.pem -subj ${CERT_SUBJ}
