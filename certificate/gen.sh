#!/bin/sh

# Remove old certificates if they exist
rm -f *.pem *.srl server-ext.cnf

# Create server-ext.cnf for X.509 v3 extensions (SANs are required for browsers)
cat > server-ext.cnf <<EOL
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = sabay.com
DNS.2 = *.sabay.com
EOL

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ca-key.pem -out ca-cert.pem \
  -subj "/C=KH/ST=Phnom Penh/L=Phnom Penh/O=Sabay/OU=Education/CN=*.sabay.com/emailAddress=mai.reaksa@sabay.com"

echo "✅ CA's self-signed certificate created"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes \
  -keyout server-key.pem -out server-req.pem \
  -subj "/C=KH/ST=Phnom Penh/L=Phnom Penh/O=PC Book/OU=Computer/CN=*.sabay.com/emailAddress=mai.reaksa@sabay.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 \
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -out server-cert.pem -extfile server-ext.cnf

echo "✅ Server's signed certificate created"
openssl x509 -in server-cert.pem -noout -text
