#!/bin/bash

# Генерация самоподписанного сертификата для HTTPS/HTTP2
echo "Generating self-signed certificate for HTTP/2..."

openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem \
  -out cert.pem \
  -days 365 \
  -nodes \
  -subj "/C=RU/ST=Moscow/L=Moscow/O=WorkMate/OU=Development/CN=localhost"

echo "Certificate generated successfully!"
echo "Files created:"
echo "  - cert.pem (certificate)"
echo "  - key.pem (private key)" 