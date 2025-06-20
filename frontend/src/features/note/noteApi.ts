import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'

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
    tagTypes: ['Note'],
    baseQuery: fetchBaseQuery({
        baseUrl: '/api',
        prepareHeaders: (headers) => {
            const token = localStorage.getItem('token')
            if (token) headers.set('Authorization', `Bearer ${token}`)
            return headers
        },
    }),
    endpoints: (builder) => ({
        getNotes: builder.query<Note[], void>({
            query: () => '/notes',
            providesTags: (result) =>
                result
                    ? [
                        ...result.map((note) => ({ type: 'Note' as const, id: note.id })),
                        { type: 'Note', id: 'LIST' },
                    ]
                    : [{ type: 'Note', id: 'LIST' }],
        }),
        getNote: builder.query<Note, string>({
            query: (id) => `/notes/${id}`,
            providesTags: (result, error, id) => [{ type: 'Note', id }],
        }),
        updateNote: builder.mutation<Note, Partial<Note> & { id: string}>({
            query: ({ id, ...patch }) => ({
                url: `/notes/${id}`,
                method: 'PUT',
                body: patch,
            }),
            invalidatesTags: (result, error, { id }) => [
                { type: 'Note', id },
                { type: 'Note', id: 'LIST' },
            ],
        }),
    }),
})

export const { useGetNotesQuery, useGetNoteQuery, useUpdateNoteMutation } = noteApi
