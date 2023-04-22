package tls

import (
   "golang.org/x/crypto/cryptobyte"
   "strings"
)

func (c *Client_Hello_Msg) UnmarshalBinary(data []byte) error {
   s := cryptobyte.String(data)
   if !s.Skip(4) { // message type and uint24 length field
      return err_fail
   }
   if !s.ReadUint16(&c.vers) { // handshake version
      return err_fail
   }
   if !s.Skip(32) { // random
      return err_fail
   }
   if !s.ReadUint8LengthPrefixed(&c.session_ID) {
      return err_fail
   }
   var cipher_suites cryptobyte.String
   if !s.ReadUint16LengthPrefixed(&cipher_suites) {
      return err_fail
   }
   for !cipher_suites.Empty() {
      var suite uint16
      if !cipher_suites.ReadUint16(&suite) {
         return err_fail
      }
      if suite == scsv_renegotiation {
         c.secure_renegotiation_supported = true
      }
      c.cipher_suites = append(c.cipher_suites, suite)
   }
   if !s.ReadUint8LengthPrefixed(&c.compression_methods) {
      return err_fail
   }
   if s.Empty() {
      // Client Hello is optionally followed by extension data
      return nil
   }
   var extensions cryptobyte.String
   if !s.ReadUint16LengthPrefixed(&extensions) || !s.Empty() {
      return err_fail
   }
   seen_exts := make(map[uint16]bool)
   for !extensions.Empty() {
      var extension uint16
      var ext_data cryptobyte.String
      if !extensions.ReadUint16(&extension) {
         return err_fail
      }
      if !extensions.ReadUint16LengthPrefixed(&ext_data) {
         return err_fail
      }
      if seen_exts[extension] {
         return err_fail
      }
      seen_exts[extension] = true
      switch extension {
      case extension_server_name:
         // RFC 6066, Section 3
         var name_list cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&name_list) || name_list.Empty() {
            return err_fail
         }
         for !name_list.Empty() {
            var name_type uint8
            var server_name cryptobyte.String
            if !name_list.ReadUint8(&name_type) {
               return err_fail
            }
            if !name_list.ReadUint16LengthPrefixed(&server_name) {
               return err_fail
            }
            if server_name.Empty() {
               return err_fail
            }
            if name_type != 0 {
               continue
            }
            if len(c.server_name) != 0 {
               // Multiple names of the same name_type are prohibited.
               return err_fail
            }
            c.server_name = string(server_name)
            // An SNI value may not include a trailing dot.
            if strings.HasSuffix(c.server_name, ".") {
               return err_fail
            }
         }
      case extension_status_request:
         // RFC 4366, Section 3.6
         var (
            ignored cryptobyte.String
            status_type uint8
         )
         if !ext_data.ReadUint8(&status_type) {
            return err_fail
         }
         if !ext_data.ReadUint16LengthPrefixed(&ignored) {
            return err_fail
         }
         if !ext_data.ReadUint16LengthPrefixed(&ignored) {
            return err_fail
         }
         c.ocsp_stapling = status_type == status_type_OCSP
      case extension_supported_curves:
         // RFC 4492, sections 5.1.1 and RFC 8446, Section 4.2.7
         var curves cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&curves) || curves.Empty() {
            return err_fail
         }
         for !curves.Empty() {
            var curve uint16
            if !curves.ReadUint16(&curve) {
               return err_fail
            }
            c.supported_curves = append(c.supported_curves, curve)
         }
      case extension_supported_points:
         // RFC 4492, Section 5.1.2
         if !ext_data.ReadUint8LengthPrefixed(&c.supported_points) {
            return err_fail
         }
         if len(c.supported_points) == 0 {
            return err_fail
         }
      case extension_session_ticket:
         // RFC 5077, Section 3.2
         c.ticket_supported = true
         ext_data.ReadBytes(&c.session_ticket, len(ext_data))
      case extension_signature_algorithms:
         // RFC 5246, Section 7.4.1.4.1
         var sig_and_algs cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&sig_and_algs) {
            return err_fail
         }
         if sig_and_algs.Empty() {
            return err_fail
         }
         for !sig_and_algs.Empty() {
            var sig_and_alg uint16
            if !sig_and_algs.ReadUint16(&sig_and_alg) {
               return err_fail
            }
            c.supported_signature_algorithms = append(
               c.supported_signature_algorithms, sig_and_alg,
            )
         }
      case extension_signature_algorithms_cert:
         // RFC 8446, Section 4.2.3
         var sig_and_algs cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&sig_and_algs) {
            return err_fail
         }
         if sig_and_algs.Empty() {
            return err_fail
         }
         for !sig_and_algs.Empty() {
            var sig_and_alg uint16
            if !sig_and_algs.ReadUint16(&sig_and_alg) {
               return err_fail
            }
            c.supported_signature_algorithms_cert = append(
               c.supported_signature_algorithms_cert, sig_and_alg,
            )
         }
      case extension_renegotiation_info:
         // RFC 5746, Section 3.2
         if !ext_data.ReadUint8LengthPrefixed(&c.secure_renegotiation) {
            return err_fail
         }
         c.secure_renegotiation_supported = true
      case extension_ALPN:
         // RFC 7301, Section 3.1
         var proto_list cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&proto_list) {
            return err_fail
         }
         if proto_list.Empty() {
            return err_fail
         }
         for !proto_list.Empty() {
            var proto cryptobyte.String
            if !proto_list.ReadUint8LengthPrefixed(&proto) || proto.Empty() {
               return err_fail
            }
            c.alpn_protocols = append(c.alpn_protocols, string(proto))
         }
      case extension_SCT:
         // RFC 6962, Section 3.3.1
         c.scts = true
      case extension_supported_versions:
         // RFC 8446, Section 4.2.1
         var vers_list cryptobyte.String
         if !ext_data.ReadUint8LengthPrefixed(&vers_list) || vers_list.Empty() {
            return err_fail
         }
         for !vers_list.Empty() {
            var vers uint16
            if !vers_list.ReadUint16(&vers) {
               return err_fail
            }
            c.supported_versions = append(c.supported_versions, vers)
         }
      case extension_cookie:
         // RFC 8446, Section 4.2.2
         if !ext_data.ReadUint16LengthPrefixed(&c.cookie) {
            return err_fail
         }
         if len(c.cookie) == 0 {
            return err_fail
         }
      case extension_key_share:
         // RFC 8446, Section 4.2.8
         var client_shares cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&client_shares) {
            return err_fail
         }
         for !client_shares.Empty() {
            var ks key_share
            if !client_shares.ReadUint16((*uint16)(&ks.group)) {
               return err_fail
            }
            if !client_shares.ReadUint16LengthPrefixed(&ks.data) {
               return err_fail
            }
            if len(ks.data) == 0 {
               return err_fail
            }
            c.key_shares = append(c.key_shares, ks)
         }
      case extension_early_data:
         // RFC 8446, Section 4.2.10
         c.early_data = true
      case extension_PSK_modes:
         // RFC 8446, Section 4.2.9
         if !ext_data.ReadUint8LengthPrefixed(&c.psk_modes) {
            return err_fail
         }
      case extension_pre_shared_key:
         // RFC 8446, Section 4.2.11
         if !extensions.Empty() {
            return err_fail // pre_shared_key must be the last extension
         }
         var identities cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&identities) {
            return err_fail
         }
         if identities.Empty() {
            return err_fail
         }
         for !identities.Empty() {
            var psk psk_identity
            if !identities.ReadUint16LengthPrefixed(&psk.label) {
               return err_fail
            }
            if !identities.ReadUint32(&psk.obfuscated_ticket_age) {
               return err_fail
            }
            if len(psk.label) == 0 {
               return err_fail
            }
            c.psk_identities = append(c.psk_identities, psk)
         }
         var binders cryptobyte.String
         if !ext_data.ReadUint16LengthPrefixed(&binders) || binders.Empty() {
            return err_fail
         }
         for !binders.Empty() {
            var binder cryptobyte.String
            if !binders.ReadUint8LengthPrefixed(&binder) {
               return err_fail
            }
            if len(binder) == 0 {
               return err_fail
            }
            c.psk_binders = append(c.psk_binders, binder)
         }
      default:
         // Ignore unknown extensions.
         continue
      }
      if !ext_data.Empty() {
         return err_fail
      }
   }
   return nil
}
