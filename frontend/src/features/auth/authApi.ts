import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

interface LoginRequest {
    email: string
    password: string
}

interface LoginResponse {
    token: string
    user: {
        id: string
        email: string
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
    }),
})

export const { useLoginMutation } = authApi
