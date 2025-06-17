import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import { useGetNotesQuery } from '@/features/note/noteApi'
import {
    Box,
    IconButton,
    Typography,
    CircularProgress,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Grid,
    Tooltip
} from '@mui/material'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward'
import NoteCard from '../features/note/components/NoteCard'
import NoteRow from '../features/note/components/NoteRow'
import React, { useState, useMemo } from 'react'
import { Link } from 'react-router-dom'

const NoteList = () => {
    const dispatch = useAppDispatch()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { data: notes, isLoading, isError } = useGetNotesQuery()

    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc') // новое состояние

    const toggleSortDirection = () => {
        setSortDirection(prev => (prev === 'asc' ? 'desc' : 'asc'))
    }

    const sortedNotes = useMemo(() => {
        if (!notes) return []

        return [...notes].sort((a, b) => {
            const aValue = new Date(a[sortBy]!).getTime()
            const bValue = new Date(b[sortBy]!).getTime()
            return sortDirection === 'asc' ? aValue - bValue : bValue - aValue
        })
    }, [notes, sortBy, sortDirection])

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box mt={4}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2} flexWrap="wrap" gap={2}>
                <Typography variant="h5" pl={2}>Мои заметки</Typography>

                <Box display="flex" alignItems="center" gap={1}>
                    <FormControl size="small">
                        <InputLabel id="sort-select-label">Сортировка</InputLabel>
                        <Select
                            labelId="sort-select-label"
                            value={sortBy}
                            label="Сортировка"
                            variant="outlined"
                            onChange={(e) => setSortBy(e.target.value as 'created_at' | 'next_review_at')}
                        >
                            <MenuItem value="created_at">По дате создания</MenuItem>
                            <MenuItem value="next_review_at">По дате следующего повторения</MenuItem>
                        </Select>
                    </FormControl>

                    <Tooltip title={`Сортировать по ${sortDirection === 'asc' ? 'возрастанию' : 'убыванию'}`}>
                        <IconButton onClick={toggleSortDirection}>
                            {sortDirection === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                        </IconButton>
                    </Tooltip>

                    <Tooltip title="Переключить вид">
                        <IconButton onClick={() => dispatch(toggleViewMode())}>
                            {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                        </IconButton>
                    </Tooltip>
                </Box>
            </Box>

            {viewMode === 'card' ? (
                <Grid mt={8} container spacing={1} justifyContent="space-between">
                    {sortedNotes.map(note => (
                        <Grid
                            key={note.id}
                            sx={{ width: { xs: '100%', sm: '100%', md: '31%' } }}
                            display="flex"
                        >
                            <Link
                                to={`/notes/${note.id}`}
                                style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}
                            >
                                <NoteCard note={note} />
                            </Link>
                        </Grid>
                    ))}
                </Grid>
            ) : (
                <Box>
                    {sortedNotes.map(note => (
                        <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                            <NoteRow note={note} />
                        </Link>
                    ))}
                </Box>
            )}
        </Box>
    )
}

export default NoteList
