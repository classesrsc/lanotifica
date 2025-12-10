package com.alessandrolattao.lanotifica.network

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

data class NotificationRequest(
    val app_name: String,
    val package_name: String,
    val title: String,
    val message: String
)

data class NotificationResponse(
    val status: String
)

interface NotificationApi {
    @POST("/notification")
    suspend fun sendNotification(@Body request: NotificationRequest): Response<NotificationResponse>
}
