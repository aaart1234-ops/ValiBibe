import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useDebounce } from 'use-debounce'
import { useGetNotesQuery } from '@/features/note/noteApi'
import { useAppSelector } from '@/app/hooks'

type SortBy = 'created_at' | 'next_review_at'
type SortDir = 'asc' | 'desc'

export function useNotesQuery() {
    const [searchParams, setSearchParams] = useSearchParams()
    const showArchived = searchParams.get('archived') === 'true'
    const pageParam = parseInt(searchParams.get('page') || '1', 10)
    const limitParam = parseInt(searchParams.get('limit') || '10', 10)

    const [sortBy, setSortBy] = useState<SortBy>('created_at')
    const [sortDirection, setSortDirection] = useState<SortDir>('desc')
    const [searchQuery, setSearchQuery] = useState('')
    const [debouncedSearchQuery] = useDebounce(searchQuery, 300)

    const [limit, setLimit] = useState(limitParam)
    const [page, setPage] = useState(pageParam - 1)
    const offset = page * limit

    const { token } = useAppSelector((s) => s.auth)

    const queryArgs = useMemo(() => ({
        search: debouncedSearchQuery,
        sortBy,
        sortDirection,
        limit,
        offset,
        archived: showArchived,
    }), [debouncedSearchQuery, sortBy, sortDirection, limit, offset, showArchived])

    const { data, isLoading, isError, refetch } = useGetNotesQuery(
        queryArgs,
        { refetchOnMountOrArgChange: true }
    )

    useEffect(() => {
        if (token) refetch()
    }, [token]) // eslint-disable-line

    useEffect(() => {
        const params: Record<string, string> = {
            page: (page + 1).toString(),
            limit: limit.toString(),
        }
        if (showArchived) params.archived = 'true'
        setSearchParams(params, { replace: false })
    }, [page, limit, showArchived]) // eslint-disable-line

    const toggleArchived = () => {
        setPage(0)
        setSearchParams((prev) => {
            const updated = new URLSearchParams(prev)
            if (showArchived) updated.delete('archived')
            else updated.set('archived', 'true')
            updated.set('page', '1')
            return updated
        })
    }

    const notes = data?.notes || []
    const total = data?.total || 0

    return {
        // data
        notes, total, isLoading, isError, refetch,
        // filters
        searchQuery, setSearchQuery,
        sortBy, setSortBy,
        sortDirection, setSortDirection,
        showArchived, toggleArchived,
        // pagination
        page, setPage,
        limit, setLimit,
    }
}
