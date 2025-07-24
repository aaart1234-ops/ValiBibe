// components/NoteList.tsx

import React, { useState, useEffect, useMemo } from 'react'
import { Link, useSearchParams, useNavigate } from 'react-router-dom'
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
    Pagination,
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
    const navigate = useNavigate()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { token } = useAppSelector((state) => state.auth)

    const [searchParams, setSearchParams] = useSearchParams()

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    // Достаём параметры из query string
    const showArchived = searchParams.get('archived') === 'true'
    const pageParam = parseInt(searchParams.get('page') || '1', 10)
    const limitParam = parseInt(searchParams.get('limit') || '10', 10)

    // Состояния фильтрации и сортировки
    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
    const [searchQuery, setSearchQuery] = useState('')
    const [debouncedSearchQuery] = useDebounce(searchQuery, 300)

    const [limit, setLimit] = useState(limitParam)
    const [page, setPage] = useState(pageParam - 1) // внутренняя нумерация с 0
    const offset = page * limit

    const toggleSortDirection = () => {
        setSortDirection((prev) => (prev === 'asc' ? 'desc' : 'asc'))
    }

    const { data, isLoading, isError, refetch } = useGetNotesQuery(
        {
            search: debouncedSearchQuery,
            sortBy,
            sortDirection,
            limit,
            offset,
            archived: showArchived,
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

    // Обновление query string при изменении page, limit, archived
    useEffect(() => {
        const params: Record<string, string> = {
            page: (page + 1).toString(),
            limit: limit.toString(),
        }
        if (showArchived) {
            params.archived = 'true'
        }
        setSearchParams(params, { replace: false })
    }, [page, limit, showArchived])

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box mt={4} pl={4} pr={4} pb={4}>
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

                    <Button
                        variant={showArchived ? 'contained' : 'outlined'}
                        color="secondary"
                        onClick={() => {
                            setPage(0)
                            setSearchParams((prev) => {
                                const updated = new URLSearchParams(prev)
                                if (showArchived) {
                                    updated.delete('archived')
                                } else {
                                    updated.set('archived', 'true')
                                }
                                updated.set('page', '1')
                                return updated
                            })
                        }}
                    >
                        {showArchived ? 'Показать активные' : 'Показать архив'}
                    </Button>

                    <IconButton onClick={() => dispatch(toggleViewMode())}>
                        {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                    </IconButton>

                    <IconButton
                        onClick={() => setOpenCreateDialog(true)}
                        color="primary"
                        size="large"
                        sx={{ padding: 2, '&:hover': { backgroundColor: '#e0e0e0' } }}
                    >
                        <AddIcon fontSize="large" />
                    </IconButton>

                    <Fab
                        color="primary"
                        onClick={() => setOpenCreateDialog(true)}
                        sx={{ position: 'fixed', bottom: 24, right: 24, zIndex: 10 }}
                    >
                        <AddIcon />
                    </Fab>
                </Box>
            </Box>

            {/* Список заметок */}
            {notes.length === 0 ? (
                <Box textAlign="center" mt={10}>
                    <Typography variant="h6" gutterBottom>
                        У вас пока нет заметок
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={2}>
                        Создайте первую заметку, чтобы начать тренироваться
                    </Typography>
                    <Button variant="contained" color="primary" component={Link} to="/notes/create">
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
                    <Select
                        value={limit}
                        onChange={(e) => {
                            const newLimit = Number(e.target.value)
                            setLimit(newLimit)
                            setPage(0)
                            setSearchParams((prev) => {
                                const updated = new URLSearchParams(prev)
                                updated.set('limit', newLimit.toString())
                                updated.set('page', '1')
                                return updated
                            })
                        }}
                    >
                        <MenuItem value={5}>5</MenuItem>
                        <MenuItem value={10}>10</MenuItem>
                        <MenuItem value={20}>20</MenuItem>
                        <MenuItem value={50}>50</MenuItem>
                    </Select>
                </FormControl>

                {total > limit && (
                    <Pagination
                        count={Math.ceil(total / limit)}
                        page={page + 1}
                        onChange={(_, value) => {
                            setPage(value - 1)
                            setSearchParams((prev) => {
                                const updated = new URLSearchParams(prev)
                                updated.set('page', value.toString())
                                return updated
                            })
                        }}
                        color="primary"
                    />
                )}
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
