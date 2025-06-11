import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export interface Note {
    id: string
    title: string
    content: string
    memoryLevel: number
    created_at: string
    next_review_at?: string
}

export const noteApi = createApi({
    reducerPath: 'noteApi',
    baseQuery: fetchBaseQuery({
        baseUrl: 'http://localhost:8081',
        prepareHeaders: (headers) => {
            const token = localStorage.getItem('token')
            if (token) headers.set('Authorization', `Bearer ${token}`)
            return headers
        },
    }),
    endpoints: (builder) => ({
        getNotes: builder.query<Note[], void>({
            query: () => '/notes',
        }),
    }),
})

export const { useGetNotesQuery } = noteApi
