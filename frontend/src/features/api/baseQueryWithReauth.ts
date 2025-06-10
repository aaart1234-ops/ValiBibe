// src/features/api/baseQueryWithReauth.ts
import { fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import type { BaseQueryFn } from '@reduxjs/toolkit/query'
import type { RootState } from '@/app/store'
import { logout } from '../auth/authSlice'

export const baseQueryWithReauth: BaseQueryFn<
    { url: string; method?: string; body?: any },
    unknown,
    unknown
> = async (args, api, extraOptions) => {
    const baseQuery = fetchBaseQuery({
        baseUrl: 'http://localhost:8081',
        prepareHeaders: (headers, { getState }) => {
            const token = (getState() as RootState).auth.token
            if (token) {
                headers.set('Authorization', `Bearer ${token}`)
            }
            return headers
        },
    })

    const result = await baseQuery(args, api, extraOptions)

    if (result.error && (result.error as any).status === 401) {
        api.dispatch(logout())
        localStorage.removeItem('token')
    }

    return result
}
