// components/NoteList.tsx

import React, { useState, useEffect } from "react"
import { Link } from 'react-router-dom'
import { useDebounce } from 'use-debounce'

import {
    Box,
    IconButton,
    Typography,
    CircularProgress,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    TextField,
    Button,
    Tooltip,
    Fab,
    Grid,
    Pagination
} from '@mui/material'

import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import AddIcon from '@mui/icons-material/Add'

import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import { useGetNotesQuery } from '@/features/note/noteApi'

import NoteCard from '@/features/note/components/NoteCard'
import NoteRow from '@/features/note/components/NoteRow'
import NoteCreateDialog from '@/features/note/components/NoteCreateDialog'

const NoteList = () => {
    const dispatch = useAppDispatch()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { token } = useAppSelector((state) => state.auth)

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    // Состояния сортировки и фильтрации
    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
    const [searchQuery, setSearchQuery] = useState('')
    const [debouncedSearchQuery] = useDebounce(searchQuery, 300)
    const [limit, setLimit] = useState(10)
    const [page, setPage] = useState(0) // 0 = первая страница
    const offset = page * limit

    const toggleSortDirection = () => {
        setSortDirection(prev => (prev === 'asc' ? 'desc' : 'asc'))
    }

    const { data, isLoading, isError, refetch } = useGetNotesQuery(
        {
            search: debouncedSearchQuery,
            sortBy,
            sortDirection,
            limit,
            offset,
        },
        {
            refetchOnMountOrArgChange: true,
        }
    )

    const notes = data?.notes || []
    const total = data?.total || 0

    useEffect(() => {
        if (token) {
            refetch()
        }
    }, [token])

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box mt={4} pl={4} pr={4} pb={4}>
            {/* Верхняя панель управления */}
            <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                mb={2}
                flexWrap="wrap"
                gap={2}
            >
                <Typography variant="h5">Мои заметки</Typography>

                <Box display="flex" alignItems="center" gap={2} flexWrap="wrap">
                    <TextField
                        label="Поиск по заголовку"
                        variant="outlined"
                        size="small"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />

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

                    <IconButton onClick={() => dispatch(toggleViewMode())}>
                        {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                    </IconButton>

                    <IconButton
                        onClick={() => setOpenCreateDialog(true)}
                        color="primary"
                        size="large"
                        sx={{
                            padding: 2,
                            '&:hover': {
                                backgroundColor: '#e0e0e0',
                            },
                        }}
                    >
                        <AddIcon fontSize="large" />
                    </IconButton>

                    <Fab
                        color="primary"
                        onClick={() => setOpenCreateDialog(true)}
                        sx={{
                            position: 'fixed',
                            bottom: 24,
                            right: 24,
                            zIndex: 10,
                        }}
                    >
                        <AddIcon />
                    </Fab>
                </Box>
            </Box>

            {/* Основной контент */}
            {notes.length === 0 ? (
                <Box textAlign="center" mt={10}>
                    <Typography variant="h6" gutterBottom>
                        У вас пока нет заметок
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={2}>
                        Создайте первую заметку, чтобы начать тренироваться
                    </Typography>
                    <Button
                        variant="contained"
                        color="primary"
                        component={Link}
                        to="/notes/create"
                    >
                        Создать заметку
                    </Button>
                </Box>
            ) : viewMode === 'card' ? (
                <Grid container spacing={2} justifyContent="start" mt={8}>
                    {notes.map((note) => (
                        <Grid
                            key={note.id}
                            sx={{ width: { xs: '100%', sm: '48%', md: '32%' } }}
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
                    {notes.map((note) => (
                        <Link
                            key={note.id}
                            to={`/notes/${note.id}`}
                            style={{ textDecoration: 'none' }}
                        >
                            <NoteRow note={note} />
                        </Link>
                    ))}
                </Box>
            )}

            {/* Пагинация */}
            <Box display="flex" justifyContent="center" alignItems="center" mt={8} gap={2}>
                <FormControl size="small">
                    <InputLabel>Показывать</InputLabel>
                    <Select
                        value={limit}
                        label="Показывать"
                        onChange={(e) => {
                            setLimit(Number(e.target.value))
                            setPage(0)
                        }}
                    >
                        <MenuItem value={5}>5</MenuItem>
                        <MenuItem value={10}>10</MenuItem>
                        <MenuItem value={20}>20</MenuItem>
                        <MenuItem value={50}>50</MenuItem>
                    </Select>
                </FormControl>

                <Pagination
                    count={Math.ceil(total / limit)}
                    page={page + 1}
                    onChange={(_, value) => setPage(value - 1)}
                    color="primary"
                />
            </Box>

            {/* Диалог создания */}
            <NoteCreateDialog
                open={openCreateDialog}
                onClose={() => setOpenCreateDialog(false)}
                onCreated={refetch}
            />
        </Box>
    )
}

export default NoteList
