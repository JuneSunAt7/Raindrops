package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/pterm/pterm"
)

func dataCert(conn net.Conn) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		pterm.Error.Println("Невозможно создать сертификат из-за", err)
		conn.Write([]byte("Невозможно создать ключ"))
		return
	}

	// Public key
	certOut, err := os.Create(CERT + "/" + Uname + "/" + Uname + ".crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()

	if err != nil {
		pterm.Warning.Println("Невозможно создать файл ", err)
	}

	defer certOut.Close()
	pterm.Success.Println(certOut.Name())

	conn.Write([]byte("1"))

	// Private key
	keyOut, err := os.OpenFile(CERT+"/"+Uname+"/"+Uname+".key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		pterm.Error.Println("FAILED CREATE CERTIFICATE FOR " + Uname)
	}
	keyOut.Close()
}

func getListCert(conn net.Conn) {
	// This function is similar to the function for getting a list of files.
	files, err := ioutil.ReadDir(CERT + "/" + Uname)
	if err != nil {
		conn.Write([]byte(err.Error()))
		pterm.Error.Println(err.Error())
		return
	}

	fileINFO := ""
	for _, file := range files {
		if !file.IsDir() {
			fileINFO += fmt.Sprintf("%-40s%-25s\n",
				file.Name(),
				file.ModTime().Format("2006-01-02 15:04:05"))
		}
	}
	conn.Write([]byte(fileINFO))

}
func CheckUserCert(uname string) bool {
	certUserPath := CERT + "/" + uname + "/" + uname + ".crt"

	if _, err := os.Stat(certUserPath); err != nil {
		if os.IsNotExist(err) {
			pterm.Error.Println("Сертификат не найден")
			return false
		} else {
			pterm.Error.Println("Сертификат не найден")
			return false
		}
	} else {
		return true
	}
}
