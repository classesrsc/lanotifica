package com.alessandrolattao.lanotifica.util

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNotNull
import org.junit.Assert.assertNull
import org.junit.Test

class QrCodeParserTest {

    @Test
    fun `parse valid QR code with token and fingerprint`() {
        val qrData = "secret123|ABC123DEF456"
        val result = QrCodeParser.parse(qrData)

        assertNotNull(result)
        assertEquals("secret123", result!!.token)
        assertEquals("ABC123DEF456", result.fingerprint)
    }

    @Test
    fun `parse returns null for single part`() {
        assertNull(QrCodeParser.parse("token"))
        assertNull(QrCodeParser.parse(""))
    }

    @Test
    fun `parse returns null for too many parts`() {
        assertNull(QrCodeParser.parse("token|fingerprint|extra"))
    }

    @Test
    fun `parse returns null for blank token`() {
        assertNull(QrCodeParser.parse("|fingerprint"))
        assertNull(QrCodeParser.parse("   |fingerprint"))
    }

    @Test
    fun `parse returns null for blank fingerprint`() {
        assertNull(QrCodeParser.parse("token|"))
        assertNull(QrCodeParser.parse("token|   "))
    }

    @Test
    fun `parse trims whitespace from all parts`() {
        val qrData = "  token123  |  FINGERPRINT  "
        val result = QrCodeParser.parse(qrData)

        assertNotNull(result)
        assertEquals("token123", result!!.token)
        assertEquals("FINGERPRINT", result.fingerprint)
    }

    @Test
    fun `parse handles real-world QR code format`() {
        val fingerprint = "A1B2C3D4E5F6A1B2C3D4E5F6A1B2C3D4E5F6A1B2C3D4E5F6A1B2C3D4E5F6A1B2"
        val token = "abc123def456abc123def456abc123def456abc123def456abc123def456abc12345"
        val qrData = "$token|$fingerprint"

        val result = QrCodeParser.parse(qrData)

        assertNotNull(result)
        assertEquals(token, result!!.token)
        assertEquals(fingerprint, result.fingerprint)
    }

    @Test
    fun `parse handles 64-char token and 64-char fingerprint`() {
        // Real format: 64-char hex token | 64-char hex SHA-256 fingerprint
        val token = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
        val fingerprint = "FEDCBA9876543210FEDCBA9876543210FEDCBA9876543210FEDCBA9876543210"
        val qrData = "$token|$fingerprint"

        val result = QrCodeParser.parse(qrData)

        assertNotNull(result)
        assertEquals(token, result!!.token)
        assertEquals(fingerprint, result.fingerprint)
    }
}
