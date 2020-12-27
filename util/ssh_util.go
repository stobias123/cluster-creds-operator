package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// func main() {
// 	savePrivateFileTo := "./id_rsa_test"
// 	savePublicFileTo := "./id_rsa_test.pub"
// 	bitSize := 4096

// 	privateKey, err := generatePrivateKey(bitSize)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

// 	err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// }
//

// GetSSHStrings returns the privatekeyPem, publicKeyBytes
func GetSSHStrings() ([]byte, []byte, error) {
	bitSize := 4096
	privateKey, err := GeneratePrivateKey(bitSize)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBytes, err := GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	privateKeyPEM := EncodePrivateKeyToPEM(privateKey)
	return privateKeyPEM, publicKeyBytes, nil
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	//log.Println("[Info] Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	//log.Println("Public key generated")
	return pubKeyBytes, nil
}