package com.alessandrolattao.lanotifica.util

import java.security.MessageDigest
import java.security.cert.X509Certificate

object CryptoUtils {
    fun calculateFingerprint(cert: X509Certificate): String {
        val md = MessageDigest.getInstance("SHA-256")
        val digest = md.digest(cert.encoded)
        return digest.joinToString("") { "%02X".format(it) }
    }

    fun calculateFingerprint(certBytes: ByteArray): String {
        val md = MessageDigest.getInstance("SHA-256")
        val digest = md.digest(certBytes)
        return digest.joinToString("") { "%02X".format(it) }
    }

    fun fingerprintsMatch(a: String, b: String): Boolean {
        return a.equals(b, ignoreCase = true)
    }
}
