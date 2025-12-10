package com.alessandrolattao.lanotifica

import android.app.Application
import android.app.NotificationChannel
import android.app.NotificationManager

class LaNotificaApp : Application() {

    companion object {
        const val CHANNEL_ID = "lanotifica_service"
    }

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
    }

    private fun createNotificationChannel() {
        val channel = NotificationChannel(
            CHANNEL_ID,
            "Notification Forwarder",
            NotificationManager.IMPORTANCE_LOW
        ).apply {
            description = "Shows when LaNotifica is forwarding notifications"
            setShowBadge(false)
        }

        val notificationManager = getSystemService(NotificationManager::class.java)
        notificationManager.createNotificationChannel(channel)
    }
}
