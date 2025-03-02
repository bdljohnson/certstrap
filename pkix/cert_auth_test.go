/*-
 * Copyright 2015 Square Inc.
 * Copyright 2014 CoreOS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pkix

import (
	"testing"
	"time"
)

func TestCreateCertificateAuthority(t *testing.T) {
	key, err := CreateRSAKey(rsaBits)
	if err != nil {
		t.Fatal("Failed creating rsa key:", err)
	}

	crt, err := CreateCertificateAuthority(key, "OU", time.Now().AddDate(5, 0, 0), "test", "US", "California", "San Francisco", "CA Name", []string{".example.com"}, nil)
	if err != nil {
		t.Fatal("Failed creating certificate authority:", err)
	}
	rawCrt, err := crt.GetRawCertificate()
	if err != nil {
		t.Fatal("Failed to get x509.Certificate:", err)
	}

	if err = rawCrt.CheckSignatureFrom(rawCrt); err != nil {
		t.Fatal("Failed to check signature:", err)
	}

	if rawCrt.Subject.OrganizationalUnit[0] != "OU" {
		t.Fatal("Failed to verify hostname:", err)
	}

	if !time.Now().After(rawCrt.NotBefore) {
		t.Fatal("Failed to be after NotBefore")
	}

	if !time.Now().Before(rawCrt.NotAfter) {
		t.Fatal("Failed to be before NotAfter")
	}

	if crt.crt.PermittedDNSDomainsCritical != true {
		t.Fatal("Permitted DNS Domains is not set to critical")
	}

	if len(crt.crt.PermittedDNSDomains) != 1 {
		t.Fatal("More than one entry found in list of permitted DNS domains")
	}

	if crt.crt.PermittedDNSDomains[0] != ".example.com" {
		t.Fatalf("Wrong permitted DNS domain, want %q, got %q", ".example.com", crt.crt.PermittedDNSDomains[0])
	}
}

func TestCreateCertificateAuthorityWithOptions(t *testing.T) {
	tests := []struct {
		name           string
		pathlen        int
		excludePathlen bool
	}{{
		name:           "pathlen: 0, excludePathlen: false",
		pathlen:        0,
		excludePathlen: false,
	}, {
		name:           "pathlen: 1, excludePathlen: false",
		pathlen:        1,
		excludePathlen: false,
	}, {
		name:           "pathlen: 0, excludePathlen: true",
		pathlen:        0,
		excludePathlen: true,
	}, {
		name:           "pathlen: 1, excludePathlen: true",
		pathlen:        1,
		excludePathlen: true,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			key, err := CreateRSAKey(rsaBits)
			if err != nil {
				t.Fatal("Failed creating rsa key:", err)
			}

			crt, err := CreateCertificateAuthorityWithOptions(key, "OU", time.Now().AddDate(5, 0, 0), "test", "US", "California", "San Francisco", "CA Name", []string{".example.com"}, nil, WithPathlenOption(tc.pathlen, tc.excludePathlen))
			if err != nil {
				t.Fatal("Failed creating certificate authority:", err)
			}
			rawCrt, err := crt.GetRawCertificate()
			if err != nil {
				t.Fatal("Failed to get x509.Certificate:", err)
			}

			if err = rawCrt.CheckSignatureFrom(rawCrt); err != nil {
				t.Fatal("Failed to check signature:", err)
			}

			if rawCrt.Subject.OrganizationalUnit[0] != "OU" {
				t.Fatal("Failed to verify hostname:", err)
			}

			if !time.Now().After(rawCrt.NotBefore) {
				t.Fatal("Failed to be after NotBefore")
			}

			if !time.Now().Before(rawCrt.NotAfter) {
				t.Fatal("Failed to be before NotAfter")
			}

			if crt.crt.PermittedDNSDomainsCritical != true {
				t.Fatal("Permitted DNS Domains is not set to critical")
			}

			if len(crt.crt.PermittedDNSDomains) != 1 {
				t.Fatal("More than one entry found in list of permitted DNS domains")
			}

			if crt.crt.PermittedDNSDomains[0] != ".example.com" {
				t.Fatalf("Wrong permitted DNS domain, want %q, got %q", ".example.com", crt.crt.PermittedDNSDomains[0])
			}

			if tc.excludePathlen {
				if crt.crt.MaxPathLen != -1 {
					t.Fatalf("Wrong MaxPathLen value, want %v, got %v", -1, crt.crt.MaxPathLen)
				}
			} else {
				if crt.crt.MaxPathLen != tc.pathlen {
					t.Fatalf("Wrong MaxPathLen value, want %v, got %v", tc.pathlen, crt.crt.MaxPathLen)
				}
			}
		})
	}
}
