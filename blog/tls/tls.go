package tls

import (
   "errors"
   "golang.org/x/crypto/cryptobyte"
)

type Client_Hello_Msg struct {
   alpn_protocols []string
   cipher_suites []uint16
   compression_methods cryptobyte.String
   cookie cryptobyte.String
   early_data bool
   key_shares []key_share
   ocsp_stapling bool
   psk_binders [][]byte
   psk_identities []psk_identity
   psk_modes cryptobyte.String
   scts bool
   secure_renegotiation cryptobyte.String
   secure_renegotiation_supported bool
   server_name string
   session_ID cryptobyte.String
   session_ticket []uint8
   supported_curves []uint16
   supported_points cryptobyte.String
   supported_signature_algorithms []uint16
   supported_signature_algorithms_cert []uint16
   supported_versions []uint16
   ticket_supported bool
   vers uint16
}

// TLS extension numbers
const (
   extension_server_name uint16 = 0
   extension_status_request uint16 = 5
   // supported_groups in TLS 1.3, see RFC 8446, Section 4.2.7
   extension_supported_curves uint16 = 10
   extension_supported_points uint16 = 11
   extension_session_ticket uint16 = 35
   extension_signature_algorithms uint16 = 13
   extension_signature_algorithms_cert uint16 = 50
   extension_renegotiation_info uint16 = 0xff01
   extension_ALPN uint16 = 16
   extension_SCT uint16 = 18
   extension_pre_shared_key uint16 = 41
   extension_early_data uint16 = 42
   extension_supported_versions uint16 = 43
   extension_cookie uint16 = 44
   extension_PSK_modes uint16 = 45
   extension_key_share uint16 = 51
)

// TLS handshake message types.
const type_client_hello uint8 = 1

var err_fail = errors.New("unmarshaling ClientHello failed")

// TLS signaling cipher suite values
const scsv_renegotiation uint16 = 0x00ff

// TLS Certificate Status Type (RFC 3546)
const status_type_OCSP uint8 = 1

func (c Client_Hello_Msg) extensions_present() bool {
   switch {
   case
   c.early_data,
   c.ocsp_stapling,
   c.scts,
   c.secure_renegotiation_supported,
   c.ticket_supported,
   len(c.alpn_protocols) >= 1,
   len(c.cookie) >= 1,
   len(c.key_shares) >= 1,
   len(c.psk_identities) >= 1,
   len(c.psk_modes) >= 1,
   len(c.server_name) >= 1,
   len(c.supported_curves) >= 1,
   len(c.supported_points) >= 1,
   len(c.supported_signature_algorithms) >= 1,
   len(c.supported_signature_algorithms_cert) >= 1,
   len(c.supported_versions) >= 1:
      return true
   }
   return false
}

// TLS 1.3 Key Share. See RFC 8446, Section 4.2.8.
type key_share struct {
   data cryptobyte.String
   group uint16
}

// TLS 1.3 PSK Identity. Can be a Session Ticket, or a reference to a saved
// session. See RFC 8446, Section 4.2.11.
type psk_identity struct {
   label cryptobyte.String
   obfuscated_ticket_age uint32
}
