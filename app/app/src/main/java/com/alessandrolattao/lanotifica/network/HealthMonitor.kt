package com.alessandrolattao.lanotifica.network

import android.content.Context
import android.util.Log
import com.alessandrolattao.lanotifica.data.SettingsRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import kotlinx.coroutines.sync.Mutex
import kotlinx.coroutines.sync.withLock
import okhttp3.OkHttpClient
import okhttp3.Request
import java.security.cert.X509Certificate
import java.util.concurrent.TimeUnit
import javax.net.ssl.SSLContext
import javax.net.ssl.TrustManager
import javax.net.ssl.X509TrustManager

/**
 * Monitors server health and manages server URL discovery.
 * Provides connection status to the UI and notification service.
 */
class HealthMonitor private constructor(private val context: Context) {

    companion object {
        private const val TAG = "HealthMonitor"
        private const val HEALTH_CHECK_INTERVAL_MS = 30_000L // 30 seconds
        private const val HEALTH_CHECK_TIMEOUT_MS = 5_000L

        @Volatile
        private var instance: HealthMonitor? = null

        fun getInstance(context: Context): HealthMonitor {
            return instance ?: synchronized(this) {
                instance ?: HealthMonitor(context.applicationContext).also { instance = it }
            }
        }
    }

    enum class ConnectionState {
        DISCONNECTED,
        CONNECTING,
        CONNECTED
    }

    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private val settingsRepository = SettingsRepository(context)
    private val serviceDiscovery = ServiceDiscovery(context)
    private val discoveryMutex = Mutex()

    private val _connectionState = MutableStateFlow(ConnectionState.DISCONNECTED)
    val connectionState: StateFlow<ConnectionState> = _connectionState.asStateFlow()

    private val _serverUrl = MutableStateFlow<String?>(null)
    val serverUrl: StateFlow<String?> = _serverUrl.asStateFlow()

    private var monitoringJob: Job? = null

    @Volatile
    private var currentFingerprint: String? = null

    /**
     * Starts monitoring server health.
     * Should be called when service is enabled.
     */
    fun startMonitoring() {
        if (monitoringJob?.isActive == true) return

        Log.d(TAG, "Starting health monitoring")
        monitoringJob = scope.launch {
            while (true) {
                checkHealth()
                delay(HEALTH_CHECK_INTERVAL_MS)
            }
        }
    }

    /**
     * Stops monitoring server health.
     */
    fun stopMonitoring() {
        Log.d(TAG, "Stopping health monitoring")
        monitoringJob?.cancel()
        monitoringJob = null
        _connectionState.value = ConnectionState.DISCONNECTED
        _serverUrl.value = null
    }

    /**
     * Forces a health check and rediscovery if needed.
     */
    suspend fun forceCheck() {
        checkHealth()
    }

    /**
     * Returns the current server URL if connected, null otherwise.
     */
    fun getServerUrlIfConnected(): String? {
        return if (_connectionState.value == ConnectionState.CONNECTED) {
            _serverUrl.value
        } else {
            null
        }
    }

    private suspend fun checkHealth() {
        val fingerprint = settingsRepository.certFingerprint.first()
        if (fingerprint.isBlank()) {
            Log.d(TAG, "No fingerprint configured, skipping health check")
            _connectionState.value = ConnectionState.DISCONNECTED
            return
        }

        currentFingerprint = fingerprint

        // Try cached URL first
        var url = _serverUrl.value ?: settingsRepository.cachedServerUrl.first().takeIf { it.isNotBlank() }

        if (url != null) {
            if (performHealthCheck(url, fingerprint)) {
                _serverUrl.value = url
                _connectionState.value = ConnectionState.CONNECTED
                return
            }
            Log.d(TAG, "Cached URL failed health check, attempting discovery")
        }

        // Discovery
        _connectionState.value = ConnectionState.CONNECTING
        url = discoverServer()

        if (url != null && performHealthCheck(url, fingerprint)) {
            _serverUrl.value = url
            settingsRepository.setCachedServerUrl(url)
            _connectionState.value = ConnectionState.CONNECTED
            Log.d(TAG, "Connected to server at $url")
        } else {
            _serverUrl.value = null
            _connectionState.value = ConnectionState.DISCONNECTED
            Log.w(TAG, "Could not connect to server")
        }
    }

    private suspend fun discoverServer(): String? {
        return discoveryMutex.withLock {
            Log.d(TAG, "Discovering server...")
            serviceDiscovery.discoverServer()?.url
        }
    }

    private fun performHealthCheck(url: String, fingerprint: String): Boolean {
        return try {
            val client = createHttpClient(fingerprint)
            val request = Request.Builder()
                .url("$url/health")
                .get()
                .build()

            val response = client.newCall(request).execute()
            val success = response.isSuccessful
            response.close()

            if (success) {
                Log.d(TAG, "Health check passed for $url")
            } else {
                Log.w(TAG, "Health check failed for $url: ${response.code}")
            }
            success
        } catch (e: Exception) {
            Log.w(TAG, "Health check error for $url: ${e.message}")
            false
        }
    }

    @Suppress("CustomX509TrustManager", "TrustAllX509TrustManager")
    private fun createHttpClient(expectedFingerprint: String): OkHttpClient {
        val trustManager = object : X509TrustManager {
            override fun checkClientTrusted(chain: Array<out X509Certificate>?, authType: String?) {}

            override fun checkServerTrusted(chain: Array<out X509Certificate>?, authType: String?) {
                if (chain.isNullOrEmpty()) {
                    throw IllegalStateException("No certificate")
                }
                val serverFingerprint = com.alessandrolattao.lanotifica.util.CryptoUtils.calculateFingerprint(chain[0])
                if (!com.alessandrolattao.lanotifica.util.CryptoUtils.fingerprintsMatch(serverFingerprint, expectedFingerprint)) {
                    throw IllegalStateException("Fingerprint mismatch")
                }
            }

            override fun getAcceptedIssuers(): Array<X509Certificate> = arrayOf()
        }

        val sslContext = SSLContext.getInstance("TLS")
        sslContext.init(null, arrayOf<TrustManager>(trustManager), null)

        return OkHttpClient.Builder()
            .sslSocketFactory(sslContext.socketFactory, trustManager)
            .hostnameVerifier { _, _ -> true }
            .connectTimeout(HEALTH_CHECK_TIMEOUT_MS, TimeUnit.MILLISECONDS)
            .readTimeout(HEALTH_CHECK_TIMEOUT_MS, TimeUnit.MILLISECONDS)
            .build()
    }
}
