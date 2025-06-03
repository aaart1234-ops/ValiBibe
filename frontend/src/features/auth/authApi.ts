import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

interface LoginRequest {
    email: string
    password: string
}

interface LoginResponse {
    token: string
}

interface RegisterRequest {
    email: string
    password: string
    nickname: string
}

interface RegisterResponse {
    message: string
    user: {
        id: string
        email: string
        nickname: string
    }
}

export const authApi = createApi({
    reducerPath: 'authApi',
    baseQuery: fetchBaseQuery({
        baseUrl: 'http://localhost:8081', // замени на твой backend URL
        credentials: 'include', // если работаешь с куками
    }),
    endpoints: (builder) => ({
        login: builder.mutation<LoginResponse, LoginRequest>({
            query: (credentials) => ({
                url: '/auth/login',
                method: 'POST',
                body: credentials,
            }),
        }),
        register: builder.mutation<RegisterResponse, RegisterRequest>({
            query: (data) => ({
                url: '/auth/register',
                method: 'POST',
                body: data,
            })
        })
    }),
})

export const { useLoginMutation, useRegisterMutation } = authApi
