package utils

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

// ParseCertificateFromPEM parse certificate from pem data
func ParseCertificateFromPEM(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(data))
	var cert *x509.Certificate
	var err error
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse cert errro")
	}
	return cert, nil
}
