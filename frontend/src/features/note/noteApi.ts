import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'

export interface Note {
    id: string
    title: string
    content: string
    memoryLevel: number
    created_at: string
    next_review_at?: string
}

export interface PaginatedNotes {
    notes: Note[]
    total: number
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
        getNotes: builder.query<PaginatedNotes, { sortBy?: string; sortDirection?: string; search?: string; limit?: number; offset?: number }>({
            query: ({ sortBy = 'created_at', sortDirection = 'desc', search = '', limit, offset } = {}) => ({
                url: '/notes',
                params: {
                    sort_by: sortBy,
                    order: sortDirection,
                    search,
                    limit,
                    offset,
                },
            }),
            providesTags: (result) =>
                result
                    ? [
                        ...result.notes.map((note: Note) => ({ type: 'Note' as const, id: note.id })),
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
        createNote: builder.mutation<Note, { title: string; content: string }>({
            query: (newNote) => ({
                url: '/notes',
                method: 'POST',
                body: newNote,
            }),
            invalidatesTags: [{ type: 'Note', id: 'LIST' }],
        }),
        deleteNote: builder.mutation<{ success: boolean }, string>({
            query: (id) => ({
                url: `/notes/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: (result, error, id) => [
                { type: 'Note', id },
                { type: 'Note', id: 'LIST' },
            ],
        }),
        archiveNote: builder.mutation<{ success: boolean }, string>({
            query: (id) => ({
                url: `/notes/${id}/archive`,
                method: 'POST',
            }),
            invalidatesTags: (result, error, id) => [
                { type: 'Note', id },
                { type: 'Note', id: 'LIST' },
            ],
        }),
    }),
})

export const {
    useGetNotesQuery,
    useGetNoteQuery,
    useUpdateNoteMutation,
    useCreateNoteMutation,
    useDeleteNoteMutation,
    useArchiveNoteMutation,
} = noteApi
