package tls

import "golang.org/x/crypto/cryptobyte"

// encoding.BinaryMarshaler
func (c *Client_Hello_Msg) MarshalBinary() ([]byte, error) {
   var b cryptobyte.Builder
   b.AddUint8(type_client_hello)
   b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
      b.AddUint16(c.vers)
      // random
      b.AddBytes(make([]byte, 32))
      b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
         b.AddBytes(c.session_ID)
      })
      b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
         for _, suite := range c.cipher_suites {
            b.AddUint16(suite)
         }
      })
      b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
         b.AddBytes(c.compression_methods)
      })
      // If extensions aren't present, omit them.
      if !c.extensions_present() {
         return
      }
      b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
         if len(c.server_name) >= 1 {
            // RFC 6066, Section 3
            b.AddUint16(extension_server_name)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  b.AddUint8(0) // name_type = host_name
                  b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                     b.AddBytes([]byte(c.server_name))
                  })
               })
            })
         }
         if c.ocsp_stapling {
            // RFC 4366, Section 3.6
            b.AddUint16(extension_status_request)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint8(1) // status_type = ocsp
               b.AddUint16(0) // empty responder_id_list
               b.AddUint16(0) // empty request_extensions
            })
         }
         if len(c.supported_curves) >= 1 {
            // RFC 4492, sections 5.1.1 and RFC 8446, Section 4.2.7
            b.AddUint16(extension_supported_curves)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, curve := range c.supported_curves {
                     b.AddUint16(uint16(curve))
                  }
               })
            })
         }
         if len(c.supported_points) >= 1 {
            // RFC 4492, Section 5.1.2
            b.AddUint16(extension_supported_points)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                  b.AddBytes(c.supported_points)
               })
            })
         }
         if c.ticket_supported {
            // RFC 5077, Section 3.2
            b.AddUint16(extension_session_ticket)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddBytes(c.session_ticket)
            })
         }
         if len(c.supported_signature_algorithms) >= 1 {
            // RFC 5246, Section 7.4.1.4.1
            b.AddUint16(extension_signature_algorithms)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, sig_algo := range c.supported_signature_algorithms {
                     b.AddUint16(uint16(sig_algo))
                  }
               })
            })
         }
         if len(c.supported_signature_algorithms_cert) >= 1 {
            // RFC 8446, Section 4.2.3
            b.AddUint16(extension_signature_algorithms_cert)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, sig_algo := range c.supported_signature_algorithms_cert {
                     b.AddUint16(uint16(sig_algo))
                  }
               })
            })
         }
         if c.secure_renegotiation_supported {
            // RFC 5746, Section 3.2
            b.AddUint16(extension_renegotiation_info)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                  b.AddBytes(c.secure_renegotiation)
               })
            })
         }
         if len(c.alpn_protocols) >= 1 {
            // RFC 7301, Section 3.1
            b.AddUint16(extension_ALPN)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, proto := range c.alpn_protocols {
                     b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                        b.AddBytes([]byte(proto))
                     })
                  }
               })
            })
         }
         if c.scts {
            // RFC 6962, Section 3.3.1
            b.AddUint16(extension_SCT)
            b.AddUint16(0) // empty extension_data
         }
         if len(c.supported_versions) >= 1 {
            // RFC 8446, Section 4.2.1
            b.AddUint16(extension_supported_versions)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, vers := range c.supported_versions {
                     b.AddUint16(vers)
                  }
               })
            })
         }
         if len(c.cookie) >= 1 {
            // RFC 8446, Section 4.2.2
            b.AddUint16(extension_cookie)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  b.AddBytes(c.cookie)
               })
            })
         }
         if len(c.key_shares) >= 1 {
            // RFC 8446, Section 4.2.8
            b.AddUint16(extension_key_share)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, ks := range c.key_shares {
                     b.AddUint16(uint16(ks.group))
                     b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                        b.AddBytes(ks.data)
                     })
                  }
               })
            })
         }
         if c.early_data {
            // RFC 8446, Section 4.2.10
            b.AddUint16(extension_early_data)
            b.AddUint16(0) // empty extension_data
         }
         if len(c.psk_modes) >= 1 {
            // RFC 8446, Section 4.2.9
            b.AddUint16(extension_PSK_modes)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                  b.AddBytes(c.psk_modes)
               })
            })
         }
         // pre_shared_key must be the last extension
         if len(c.psk_identities) >= 1 {
            // RFC 8446, Section 4.2.11
            b.AddUint16(extension_pre_shared_key)
            b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, psk := range c.psk_identities {
                     b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                        b.AddBytes(psk.label)
                     })
                     b.AddUint32(psk.obfuscated_ticket_age)
                  }
               })
               b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
                  for _, binder := range c.psk_binders {
                     b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
                        b.AddBytes(binder)
                     })
                  }
               })
            })
         }
      })
   })
   return b.Bytes()
}
