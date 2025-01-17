package print

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/effective-security/xpki/certutil"
	"github.com/effective-security/xpki/oid"
	"golang.org/x/crypto/ocsp"
)

// Certificates prints list of cert details
func Certificates(w io.Writer, list []*x509.Certificate, verboseExtensions bool) {
	for idx, crt := range list {
		fmt.Fprintf(w, "==================================== %d ====================================\n", 1+idx)
		Certificate(w, crt, verboseExtensions)
	}
}

// Certificate prints cert details
func Certificate(w io.Writer, crt *x509.Certificate, verboseExtensions bool) {
	now := time.Now()
	issuedIn := now.Sub(crt.NotBefore.Local()) / time.Minute * time.Minute
	expiresIn := crt.NotAfter.Local().Sub(now) / time.Minute * time.Minute

	fmt.Fprintf(w, "SKID: %s\n", certutil.GetSubjectKeyID(crt))
	fmt.Fprintf(w, "IKID: %s\n", certutil.GetAuthorityKeyID(crt))
	fmt.Fprintf(w, "Subject: %s\n", certutil.NameToString(&crt.Subject))
	fmt.Fprintf(w, "Serial: %s\n", crt.SerialNumber.String())
	fmt.Fprintf(w, "Issuer: %s\n", certutil.NameToString(&crt.Issuer))
	fmt.Fprintf(w, "Issued: %s (%s ago)\n", crt.NotBefore.Local().String(), issuedIn.String())
	fmt.Fprintf(w, "Expires: %s (in %s)\n", crt.NotAfter.Local().String(), expiresIn.String())

	if len(crt.DNSNames) > 0 {
		fmt.Fprintf(w, "DNS Names:\n")
		for _, n := range crt.DNSNames {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.IPAddresses) > 0 {
		fmt.Fprintf(w, "IP Addresses:\n")
		for _, n := range crt.IPAddresses {
			fmt.Fprintf(w, "  - %s\n", n.String())
		}
	}
	if len(crt.URIs) > 0 {
		fmt.Fprintf(w, "URIs:\n")
		for _, n := range crt.URIs {
			fmt.Fprintf(w, "  - %s\n", n.String())
		}
	}
	if len(crt.EmailAddresses) > 0 {
		fmt.Fprintf(w, "Emails:\n")
		for _, n := range crt.EmailAddresses {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.CRLDistributionPoints) > 0 {
		fmt.Fprintf(w, "CRL Distribution Points:\n")
		for _, n := range crt.CRLDistributionPoints {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.OCSPServer) > 0 {
		fmt.Fprintf(w, "OCSP Servers:\n")
		for _, n := range crt.OCSPServer {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.IssuingCertificateURL) > 0 {
		fmt.Fprintf(w, "Issuing Certificates:\n")
		for _, n := range crt.IssuingCertificateURL {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	fmt.Fprintf(w, "CA: %t\n", crt.IsCA)
	fmt.Fprintf(w, "  Basic Constraints Valid: %t\n", crt.BasicConstraintsValid)
	fmt.Fprintf(w, "  Max Path: %d\n", crt.MaxPathLen)

	if verboseExtensions && len(crt.Extensions) > 0 {
		fmt.Fprintf(w, "Extensions:\n")
		for _, ex := range crt.Extensions {
			soid := ex.Id.String()
			fmt.Fprintf(w, "  critical: %t\n", ex.Critical)
			if name, ok := oid.DisplayName[soid]; ok {
				fmt.Fprintf(w, "  oid: %s (%s)\n", soid, name)
			} else {
				fmt.Fprintf(w, "  oid: %s\n", soid)
			}

			fmt.Fprintf(w, "  value: %x\n", ex.Value)
			switch soid {
			case "2.5.29.15":
				fmt.Fprintf(w, "  - %s\n", strings.Join(oid.KeyUsages(crt.KeyUsage), ", "))
			case "2.5.29.37":
				fmt.Fprintf(w, "  - %s\n", strings.Join(oid.ExtKeyUsages(crt.ExtKeyUsage...), ", "))
			case "2.5.29.32":
				fmt.Fprintf(w, "  identifiers: %s\n", strings.Join(oid.Strings(crt.PolicyIdentifiers...), ", "))
			}
			fmt.Fprintln(w)
		}
	}
}

// CertificateRequest prints cert request details
func CertificateRequest(w io.Writer, crt *x509.CertificateRequest) {
	fmt.Fprintf(w, "Subject: %s\n", certutil.NameToString(&crt.Subject))
	if len(crt.DNSNames) > 0 {
		fmt.Fprintf(w, "DNS Names:\n")
		for _, n := range crt.DNSNames {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.IPAddresses) > 0 {
		fmt.Fprintf(w, "IP Addresses:\n")
		for _, n := range crt.IPAddresses {
			fmt.Fprintf(w, "  - %s\n", n.String())
		}
	}
	if len(crt.URIs) > 0 {
		fmt.Fprintf(w, "URIs:\n")
		for _, n := range crt.URIs {
			fmt.Fprintf(w, "  - %s\n", n.String())
		}
	}
	if len(crt.EmailAddresses) > 0 {
		fmt.Fprintf(w, "Emails:\n")
		for _, n := range crt.EmailAddresses {
			fmt.Fprintf(w, "  - %s\n", n)
		}
	}
	if len(crt.Extensions) > 0 {
		fmt.Fprintf(w, "Extensions:\n")
		for _, n := range crt.Extensions {
			fmt.Fprintf(w, "  - %s\n", n.Id.String())
		}
	}
}

// CertificateList prints CRL details
func CertificateList(w io.Writer, crl *pkix.CertificateList) {
	now := time.Now()
	issuedIn := now.Sub(crl.TBSCertList.ThisUpdate) / time.Minute * time.Minute
	expiresIn := crl.TBSCertList.NextUpdate.Sub(now) / time.Minute * time.Minute

	fmt.Fprintf(w, "Version: %d\n", crl.TBSCertList.Version)
	fmt.Fprintf(w, "Issuer: %s\n", crl.TBSCertList.Issuer.String())
	fmt.Fprintf(w, "Issued: %s (%s ago)\n", crl.TBSCertList.ThisUpdate.Local().String(), issuedIn.String())
	fmt.Fprintf(w, "Expires: %s (in %s)\n", crl.TBSCertList.NextUpdate.Local().String(), expiresIn.String())

	if len(crl.TBSCertList.RevokedCertificates) > 0 {
		fmt.Fprintf(w, "Revoked:\n")
		for _, r := range crl.TBSCertList.RevokedCertificates {
			fmt.Fprintf(w, "  - %s | %s\n",
				r.SerialNumber.String(),
				r.RevocationTime.Local().Format(time.RFC3339))
		}
	}
}

var ocspStatusCode = map[int]string{
	ocsp.Good:    "good",
	ocsp.Revoked: "revoked",
	ocsp.Unknown: "unknown",
}

// OCSPResponse prints OCSP response details
func OCSPResponse(w io.Writer, res *ocsp.Response, verboseExtensions bool) {
	now := time.Now()
	issuedIn := now.Sub(res.ProducedAt) / time.Minute * time.Minute
	updatedIn := now.Sub(res.ThisUpdate) / time.Minute * time.Minute
	expiresIn := res.NextUpdate.Sub(now) / time.Minute * time.Minute

	fmt.Fprintf(w, "Serial: %s\n", res.SerialNumber.String())
	fmt.Fprintf(w, "Issued: %s (%s ago)\n", res.ProducedAt.Local().String(), issuedIn.String())
	fmt.Fprintf(w, "Updated: %s (%s ago)\n", res.ThisUpdate.Local().String(), updatedIn.String())
	fmt.Fprintf(w, "Expires: %s (in %s)\n", res.NextUpdate.Local().String(), expiresIn.String())
	fmt.Fprintf(w, "Status: %s\n", ocspStatusCode[res.Status])
	if res.Status == ocsp.Revoked {
		fmt.Fprintf(w, "Revocation reason: %d\n", res.RevocationReason)
		revokedIn := now.Sub(res.RevokedAt) / time.Minute * time.Minute
		fmt.Fprintf(w, "Revoked: %s (%s ago)\n", res.RevokedAt.Local().String(), revokedIn.String())
	}
	if verboseExtensions {
		if len(res.RawResponderName) > 0 {
			fmt.Fprintf(w, "Responder name hash: %x\n", res.RawResponderName)
		}
		if len(res.ResponderKeyHash) > 0 {
			fmt.Fprintf(w, "Responder key hash: %x\n", res.ResponderKeyHash)
		}
	}
	if verboseExtensions && len(res.Extensions) > 0 {
		fmt.Fprintf(w, "Extensions:\n")
		for _, ex := range res.Extensions {
			fmt.Fprintf(w, "  id: %s, critical: %t\n", ex.Id, ex.Critical)
			fmt.Fprintf(w, "  val: %x\n\n", ex.Value)
		}
	}
	if res.Certificate != nil {
		fmt.Fprintf(w, "Certificate:\n")
		Certificate(w, res.Certificate, verboseExtensions)
	}
}

// CSRandCert outputs a cert, key and csr
func CSRandCert(w io.Writer, key, csrBytes, cert []byte) {
	out := map[string]string{}
	if cert != nil {
		out["cert"] = string(cert)
	}

	if key != nil {
		out["key"] = string(key)
	}

	if csrBytes != nil {
		out["csr"] = string(csrBytes)
	}

	jsonOut, _ := json.Marshal(out)
	fmt.Fprintln(w, string(jsonOut))
}
