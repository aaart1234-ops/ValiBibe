// components/NoteList.tsx
import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import { useGetNotesQuery } from '@/features/note/noteApi'
import { Box, IconButton, Typography, CircularProgress } from '@mui/material'
import Grid from '@mui/material/Grid'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import NoteCard from '../features/note/components/NoteCard'
import NoteRow from '../features/note/components/NoteRow'
import React from "react";
import { Link } from 'react-router-dom'

const NoteList = () => {
    const dispatch = useAppDispatch()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { data: notes, isLoading, isError } = useGetNotesQuery()

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h5">Мои заметки</Typography>
                <IconButton onClick={() => dispatch(toggleViewMode())}>
                    {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                </IconButton>
            </Box>

            {viewMode === 'card' ? (
                <Grid container spacing={2}>
                    {notes.map(note => (
                        <Grid key={note.id} sx={{ width: { xs: '100%', sm: '48%', md: '33.33%' } }}>
                            <Link to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                                <NoteCard note={note} />
                            </Link>
                        </Grid>
                    ))}
                </Grid>
            ) : (
                <Box>
                    {notes.map(note => (
                        <Link to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                            <NoteRow key={note.id} note={note} />
                        </Link>
                    ))}
                </Box>
            )}
        </Box>
    )
}

export default NoteList
