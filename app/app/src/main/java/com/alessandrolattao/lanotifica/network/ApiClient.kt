package com.alessandrolattao.lanotifica.network

import android.annotation.SuppressLint
import android.util.Log
import com.alessandrolattao.lanotifica.util.CryptoUtils
import com.alessandrolattao.lanotifica.util.UrlUtils
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.security.cert.X509Certificate
import java.util.concurrent.TimeUnit
import javax.net.ssl.SSLContext
import javax.net.ssl.TrustManager
import javax.net.ssl.X509TrustManager

object ApiClient {

    private const val TAG = "ApiClient"

    @Volatile
    private var currentBaseUrl: String? = null

    @Volatile
    private var currentToken: String? = null

    @Volatile
    private var currentFingerprint: String? = null

    @Volatile
    private var api: NotificationApi? = null

    fun getApi(baseUrl: String, token: String, expectedFingerprint: String): NotificationApi {
        val normalizedUrl = UrlUtils.normalizeUrl(baseUrl)

        if (api != null &&
            currentBaseUrl == normalizedUrl &&
            currentToken == token &&
            currentFingerprint == expectedFingerprint) {
            return api!!
        }

        synchronized(this) {
            val client = createClientWithPinning(expectedFingerprint, token)

            val retrofit = Retrofit.Builder()
                .baseUrl(normalizedUrl)
                .client(client)
                .addConverterFactory(GsonConverterFactory.create())
                .build()

            currentBaseUrl = normalizedUrl
            currentToken = token
            currentFingerprint = expectedFingerprint
            api = retrofit.create(NotificationApi::class.java)
            return api!!
        }
    }

    @SuppressLint("CustomX509TrustManager", "TrustAllX509TrustManager")
    private fun createClientWithPinning(expectedFingerprint: String, token: String): OkHttpClient {
        val trustManager = object : X509TrustManager {
            override fun checkClientTrusted(chain: Array<out X509Certificate>?, authType: String?) {}

            override fun checkServerTrusted(chain: Array<out X509Certificate>?, authType: String?) {
                if (chain.isNullOrEmpty()) {
                    throw IllegalStateException("No certificate provided by server")
                }

                val serverCert = chain[0]
                val serverFingerprint = CryptoUtils.calculateFingerprint(serverCert)

                Log.d(TAG, "Server certificate fingerprint: $serverFingerprint")
                Log.d(TAG, "Expected fingerprint: $expectedFingerprint")

                if (!CryptoUtils.fingerprintsMatch(serverFingerprint, expectedFingerprint)) {
                    throw IllegalStateException(
                        "Certificate fingerprint mismatch! Expected: $expectedFingerprint, Got: $serverFingerprint"
                    )
                }

                Log.d(TAG, "Certificate pinning verification passed")
            }

            override fun getAcceptedIssuers(): Array<X509Certificate> = arrayOf()
        }

        val sslContext = SSLContext.getInstance("TLS")
        sslContext.init(null, arrayOf<TrustManager>(trustManager), null)

        return OkHttpClient.Builder()
            .sslSocketFactory(sslContext.socketFactory, trustManager)
            .hostnameVerifier { _, _ -> true } // We verify via fingerprint instead
            .connectTimeout(10, TimeUnit.SECONDS)
            .readTimeout(10, TimeUnit.SECONDS)
            .writeTimeout(10, TimeUnit.SECONDS)
            .addInterceptor(authInterceptor(token))
            .addInterceptor(
                HttpLoggingInterceptor().apply {
                    level = HttpLoggingInterceptor.Level.BODY
                }
            )
            .build()
    }

    private fun authInterceptor(token: String): Interceptor {
        return Interceptor { chain ->
            val request = chain.request().newBuilder()
                .addHeader("Authorization", "Bearer $token")
                .build()
            chain.proceed(request)
        }
    }
}
