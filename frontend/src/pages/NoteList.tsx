// pages/NoteList.tsx
import React, { useEffect, useState } from 'react'
import { Box, CircularProgress, Typography } from '@mui/material'
import { useTheme } from '@mui/material/styles'
import useMediaQuery from '@mui/material/useMediaQuery'
import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import { useSearchParams } from 'react-router-dom'

import { useNotesQuery } from '@/features/note/hooks/useNotesQuery'
import { useArchiveWithUndo } from '@/features/note/hooks/useArchiveWithUndo'

import NoteFilters from '@/features/note/components/NoteFilters'
import NotesView from '@/features/note/components/NotesView'
import NotePagination from '@/features/note/components/NotePagination'
import UndoSnackbar from '@/features/note/components/UndoSnackbar'
import NoteCreateDialog from '@/features/note/components/NoteCreateDialog'
import NoteDeleteDialog from '@/features/note/components/NoteDeleteDialog'

import { Note } from '@/features/note/noteApi'

interface NoteListProps {
    isArchiveView?: boolean
}

const NoteList: React.FC<NoteListProps> = ({ isArchiveView = false }) => {
    const dispatch = useAppDispatch()
    const viewMode = useAppSelector((s) => s.notes.viewMode) as 'card' | 'list'
    const [searchParams, setSearchParams] = useSearchParams()

    // синхронизация isArchiveView с URL параметром archived
    useEffect(() => {
        setSearchParams((prev) => {
            const updated = new URLSearchParams(prev)
            if (isArchiveView) {
                updated.set('archived', 'true')
                updated.set('page', '1')
            } else {
                updated.delete('archived')
            }
            return updated
        })
    }, [isArchiveView, setSearchParams])

    const {
        notes,
        total,
        isLoading,
        isError,
        refetch,
        searchQuery,
        setSearchQuery,
        sortBy,
        setSortBy,
        sortDirection,
        setSortDirection,
        page,
        setPage,
        limit,
        setLimit,
    } = useNotesQuery()

    const {
        optimisticArchivedIds,
        snackbarOpen,
        handleSnackbarClose,
        handleRequestArchive,
        handleUndo,
        autoHideDuration,
    } = useArchiveWithUndo({ onCommitted: refetch })

    const theme = useTheme()
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    // Диалог удаления
    const [noteToDelete, setNoteToDelete] = useState<Note | null>(null)

    if (isLoading) return <CircularProgress />
    if (isError) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box pb={4}>
            <NoteFilters
                searchQuery={searchQuery}
                onSearchChange={setSearchQuery}
                sortBy={sortBy}
                onSortByChange={setSortBy}
                sortDirection={sortDirection}
                onToggleSortDirection={() => setSortDirection((p) => (p === 'asc' ? 'desc' : 'asc'))}
                viewMode={viewMode}
                onToggleViewMode={() => dispatch(toggleViewMode())}
                onOpenCreateDialog={() => setOpenCreateDialog(true)}
            />

            <NotesView
                notes={notes.filter((n) => !optimisticArchivedIds.has(n.id))}
                viewMode={viewMode}
                isMobile={isMobile}
                onRequestArchive={handleRequestArchive}
                onRequestDelete={(note) => setNoteToDelete(note)}
                onRefetch={refetch}
                isArchiveView={isArchiveView}
            />

            <NotePagination
                total={total}
                limit={limit}
                page={page}
                onLimitChange={(newLimit) => {
                    setLimit(newLimit)
                    setPage(0)
                }}
                onPageChange={(newPage) => setPage(newPage)}
            />

            <NoteCreateDialog
                open={openCreateDialog}
                onClose={() => setOpenCreateDialog(false)}
                onCreated={refetch}
            />

            <UndoSnackbar
                open={snackbarOpen}
                autoHideDuration={autoHideDuration}
                onClose={handleSnackbarClose}
                onUndo={handleUndo}
            />

            <NoteDeleteDialog
                note={noteToDelete}
                open={!!noteToDelete}
                onClose={() => setNoteToDelete(null)}
                onDeleted={refetch}
            />
        </Box>
    )
}

export default NoteList
