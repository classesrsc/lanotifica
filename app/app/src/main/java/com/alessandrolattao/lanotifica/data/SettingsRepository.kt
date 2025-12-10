package com.alessandrolattao.lanotifica.data

import android.content.Context
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

private val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")

class SettingsRepository(private val context: Context) {

    companion object {
        private val AUTH_TOKEN = stringPreferencesKey("auth_token")
        private val CERT_FINGERPRINT = stringPreferencesKey("cert_fingerprint")
        private val SERVICE_ENABLED = booleanPreferencesKey("service_enabled")
        private val CACHED_SERVER_URL = stringPreferencesKey("cached_server_url")
    }

    val authToken: Flow<String> = context.dataStore.data.map { preferences ->
        preferences[AUTH_TOKEN] ?: ""
    }

    val certFingerprint: Flow<String> = context.dataStore.data.map { preferences ->
        preferences[CERT_FINGERPRINT] ?: ""
    }

    val serviceEnabled: Flow<Boolean> = context.dataStore.data.map { preferences ->
        preferences[SERVICE_ENABLED] ?: false
    }

    val cachedServerUrl: Flow<String> = context.dataStore.data.map { preferences ->
        preferences[CACHED_SERVER_URL] ?: ""
    }

    val isConfigured: Flow<Boolean> = context.dataStore.data.map { preferences ->
        !preferences[AUTH_TOKEN].isNullOrBlank() && !preferences[CERT_FINGERPRINT].isNullOrBlank()
    }

    suspend fun setServerConfig(token: String, fingerprint: String) {
        context.dataStore.edit { preferences ->
            preferences[AUTH_TOKEN] = token
            preferences[CERT_FINGERPRINT] = fingerprint
            preferences.remove(CACHED_SERVER_URL) // Clear cached URL to force rediscovery
        }
    }

    suspend fun setCachedServerUrl(url: String) {
        context.dataStore.edit { preferences ->
            preferences[CACHED_SERVER_URL] = url
        }
    }

    suspend fun setServiceEnabled(enabled: Boolean) {
        context.dataStore.edit { preferences ->
            preferences[SERVICE_ENABLED] = enabled
        }
    }
}
