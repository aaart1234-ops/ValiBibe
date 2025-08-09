// components/NoteList.tsx
import React, { useEffect, useRef, useState } from 'react'
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
    Snackbar,
    Alert,
} from '@mui/material'

import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import AddIcon from '@mui/icons-material/Add'
import ArchiveIcon from '@mui/icons-material/Inventory2'
import ArchiveOutlinedIcon from '@mui/icons-material/Inventory2Outlined'

import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import {
    useGetNotesQuery,
    useArchiveNoteMutation,
    useUnarchiveNoteMutation,
    Note,
} from '@/features/note/noteApi'

import NoteCard from '@/features/note/components/NoteCard'
import NoteRow from '@/features/note/components/NoteRow'
import SwipeableNoteCard from '@/features/note/components/SwipeableNoteCard'
import SwipeableNoteRow from '@/features/note/components/SwipeableNoteRow'
import NoteCreateDialog from '@/features/note/components/NoteCreateDialog'
import useMediaQuery from '@mui/material/useMediaQuery'
import { useTheme } from '@mui/material/styles'

const ARCHIVE_DELAY_MS = 4000 // сколько ждать до фактического архивирования (даём время на Undo)

const NoteList: React.FC = () => {
    const dispatch = useAppDispatch()
    const navigate = useNavigate()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { token } = useAppSelector((state) => state.auth)

    const [searchParams, setSearchParams] = useSearchParams()
    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    // Query params
    const showArchived = searchParams.get('archived') === 'true'
    const pageParam = parseInt(searchParams.get('page') || '1', 10)
    const limitParam = parseInt(searchParams.get('limit') || '10', 10)

    // Sorting / filtering / search
    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
    const [searchQuery, setSearchQuery] = useState('')
    const [debouncedSearchQuery] = useDebounce(searchQuery, 300)

    // Pagination
    const [limit, setLimit] = useState(limitParam)
    const [page, setPage] = useState(pageParam - 1)
    const offset = page * limit

    // Archive / Snackbar / Undo state (moved here)
    const [archiveNote] = useArchiveNoteMutation()
    const [unarchiveNote] = useUnarchiveNoteMutation()

    // pendingArchivedIds: скрываем заметки из UI пока ожидается завершение архивирования
    const [pendingArchivedIds, setPendingArchivedIds] = useState<string[]>([])
    // lastArchivedNote — заметка, для которой показывается Snackbar (Undo)
    const [lastArchivedNote, setLastArchivedNote] = useState<Note | null>(null)
    const [snackbarOpen, setSnackbarOpen] = useState(false)
    // timersRef хранит отложенные таймеры по id заметки
    const timersRef = useRef<Record<string, ReturnType<typeof setTimeout> | undefined>>({})

    const toggleSortDirection = () => {
        setSortDirection((prev) => (prev === 'asc' ? 'desc' : 'asc'))
    }

    const handleToggleArchived = () => {
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
        { refetchOnMountOrArgChange: true }
    )

    const notes = data?.notes || []
    const total = data?.total || 0

    const theme = useTheme()
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'))

    useEffect(() => {
        if (token) {
            refetch()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [token])

    // Cleanup таймеров при размонтировании
    useEffect(() => {
        return () => {
            Object.values(timersRef.current).forEach((t) => {
                if (t) clearTimeout(t)
            })
            timersRef.current = {}
        }
    }, [])

    // Обновление query string при изменении page/limit/archived
    useEffect(() => {
        const params: Record<string, string> = {
            page: (page + 1).toString(),
            limit: limit.toString(),
        }
        if (showArchived) params.archived = 'true'
        setSearchParams(params, { replace: false })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [page, limit, showArchived])

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    // --- Archive / Undo handlers (в NoteList, с Snackbar) ---

    // Запрос архивирования — вызывается из Swipeable- компонентов
    const handleRequestArchive = (note: Note) => {
        // Скрываем заметку из UI сразу (пока ожидаем подтверждение/timeout)
        setPendingArchivedIds((prev) => (prev.includes(note.id) ? prev : [...prev, note.id]))

        // Показываем snackbar для этой заметки
        setLastArchivedNote(note)
        setSnackbarOpen(true)

        // Создаём таймер — через ARCHIVE_DELAY_MS выполняем реальную мутацию
        const timer = setTimeout(async () => {
            try {
                await archiveNote(note.id).unwrap()
                // перезапросим список
                refetch()
            } catch (e) {
                console.error('Ошибка при архивировании:', e)
            } finally {
                // убираем pending id и очистим таймер
                setPendingArchivedIds((prev) => prev.filter((id) => id !== note.id))
                delete timersRef.current[note.id]
                // если это была последняя отображаемая snackbar-заметка — очищаем
                setLastArchivedNote((curr) => (curr?.id === note.id ? null : curr))
            }
        }, ARCHIVE_DELAY_MS)

        timersRef.current[note.id] = timer
    }

    // Если нужно срочно разархивировать (например свайп вправо по уже заархивированной заметке)
    const handleRequestUnarchive = async (note: Note) => {
        try {
            await unarchiveNote(note.id).unwrap()
            await refetch()
        } catch (e) {
            console.error('Ошибка при разархивировании:', e)
        }
    }

    // Undo — отменяем отложенное архивирование (если таймер ещё не сработал),
    // либо если архив уже применён — делаем разархивирование на сервере.
    const handleUndo = async () => {
        const note = lastArchivedNote
        if (!note) {
            setSnackbarOpen(false)
            return
        }

        const timer = timersRef.current[note.id]

        if (timer) {
            // таймер ещё не сработал — отменяем архивирование
            clearTimeout(timer)
            delete timersRef.current[note.id]
            setPendingArchivedIds((prev) => prev.filter((id) => id !== note.id))
            setLastArchivedNote(null)
            setSnackbarOpen(false)
            // не нужно вызывать API — мы отменили отложенную мутацию
            return
        }

        // таймер уже сработал и заметка, возможно, уже заархивирована на сервере
        try {
            await unarchiveNote(note.id).unwrap()
            await refetch()
        } catch (e) {
            console.error('Ошибка при отмене архивирования (unarchive):', e)
        } finally {
            setLastArchivedNote(null)
            setSnackbarOpen(false)
        }
    }

    const handleSnackbarClose = (_event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') return
        setSnackbarOpen(false)
    }

    // --- Рендер ---
    // Скрываем заметки, которые находятся в pendingArchivedIds (они визуально "уходят" до завершения архивации)
    const visibleNotes = notes.filter((n) => !pendingArchivedIds.includes(n.id))

    return (
        <Box
            mt={4}
            pb={4}
            sx={{
                pl: isMobile ? 1 : 4,
                pr: isMobile ? 1 : 4,
            }}
        >
            <Box display="flex" flexDirection="column" justifyContent="space-between" alignItems="space-between" gap={2} mb={2}>
                {/* Заголовок */}
                <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="h5">Мои заметки</Typography>
                </Box>

                <Box
                    display="flex"
                    flexDirection={isMobile ? 'column' : 'row'}
                    alignItems="space-between"
                    justifyContent="space-between"
                    gap={0}
                    flexWrap="wrap"
                >
                    {/* Поиск + сортировка */}
                    <Box
                        display="flex"
                        flexDirection="row"
                        alignItems="space-between"
                        justifyContent="space-between"
                        gap={2}
                        flexWrap="wrap"
                        flexGrow={1}
                        minWidth={isMobile ? '100%' : 'auto'}
                    >
                        <TextField
                            label="Поиск по заголовку"
                            variant="outlined"
                            size="small"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            sx={{ width: isMobile ? '100%' : 280 }}
                        />

                        <FormControl size="small" sx={{ minWidth: 160 }}>
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

                        <Tooltip
                            title={`Сортировать по ${sortDirection === 'asc' ? 'возрастанию' : 'убыванию'}`}
                            sx={{ backgroundColor: '#e0e0e0' }}
                            color="primary"
                        >
                            <IconButton onClick={toggleSortDirection}>
                                {sortDirection === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                            </IconButton>
                        </Tooltip>

                        <Tooltip
                            title={viewMode === 'card' ? 'Список' : 'Карточки'}
                            sx={{ backgroundColor: '#e0e0e0' }}
                            color="primary"
                        >
                            <IconButton
                                onClick={() => dispatch(toggleViewMode())}
                            >
                                {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                            </IconButton>
                        </Tooltip>
                        {/* </Box>*/}

                        {/* Правая часть: архив, вид, добавить */}
                        {/*<Box display="flex" alignItems="center" gap={1} flexWrap="wrap">*/}

                        {isMobile ? (
                            <Tooltip
                                title={showArchived ? 'Показать активные' : 'Показать архив'}
                                sx={{ backgroundColor: '#e0e0e0' }}
                            >
                                <IconButton color="secondary" onClick={handleToggleArchived}>
                                    {showArchived ? <ArchiveOutlinedIcon /> : <ArchiveIcon />}
                                </IconButton>
                            </Tooltip>
                        ) : (
                            <Button
                                variant={showArchived ? 'contained' : 'outlined'}
                                color="primary"
                                onClick={handleToggleArchived}
                            >
                                {showArchived ? 'Показать активные' : 'Показать архив'}
                            </Button>
                        )}

                        {!isMobile && (
                            <IconButton
                                onClick={() => setOpenCreateDialog(true)}
                                color="primary"
                                size="large"
                                sx={{ padding: 0, '&:hover': { backgroundColor: '#e0e0e0' } }}
                            >
                                <AddIcon fontSize="large" />
                            </IconButton>
                        )}
                    </Box>
                </Box>

                {/* FAB на мобильных */}
                {isMobile && (
                    <Fab color="primary" onClick={() => setOpenCreateDialog(true)} sx={{ position: 'fixed', bottom: 24, right: 24, zIndex: 10 }}>
                        <AddIcon />
                    </Fab>
                )}
            </Box>

            {/* Список заметок */}
            {visibleNotes.length === 0 ? (
                <Box textAlign="center" mt={10}>
                    <Typography variant="h6" gutterBottom>
                        У вас пока нет заметок
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={2}>
                        Создайте первую заметку, чтобы начать тренироваться
                    </Typography>
                    <Button variant="contained" color="primary" component={Link} to="/notes/create">Создать заметку</Button>
                </Box>
            ) : viewMode === 'card' ? (
                <Grid container spacing={2} justifyContent="start" mt={4}>
                    {visibleNotes.map((note) => (
                        <Grid key={note.id} sx={{ width: { xs: '100%', sm: '48%', md: '32%' } }} display="flex">
                            {isMobile ? (
                                // Swipeable компонент теперь должен вызвать onArchiveRequest / onUnarchiveRequest
                                <SwipeableNoteCard
                                    note={note}
                                    onRefetch={refetch}
                                    onArchiveRequest={handleRequestArchive}
                                    onUnarchiveRequest={handleRequestUnarchive}
                                />
                            ) : (
                                <Link to={`/notes/${note.id}`} style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}>
                                    <NoteCard note={note} />
                                </Link>
                            )}
                        </Grid>
                    ))}
                </Grid>
            ) : (
                <Box mt={4}>
                    {visibleNotes.map((note) =>
                        isMobile ? (
                            <SwipeableNoteRow
                                key={note.id}
                                note={note}
                                onRefetch={refetch}
                                onArchiveRequest={handleRequestArchive}
                                onUnarchiveRequest={handleRequestUnarchive}
                            />
                        ) : (
                            <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                                <NoteRow note={note} />
                            </Link>
                        )
                    )}
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
            <NoteCreateDialog open={openCreateDialog} onClose={() => setOpenCreateDialog(false)} onCreated={refetch} />

            {/* Snackbar для архивирования (Undo) */}
            <Snackbar open={snackbarOpen} autoHideDuration={ARCHIVE_DELAY_MS} onClose={handleSnackbarClose} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                <Alert
                    onClose={handleSnackbarClose}
                    severity="info"
                    sx={{ width: '100%' }}
                    action={
                        <Button color="inherit" size="small" onClick={handleUndo}>
                            Отменить
                        </Button>
                    }
                >
                    {lastArchivedNote ? `Заметка "${lastArchivedNote.title}" перемещена в архив` : 'Заметка перемещена в архив'}
                </Alert>
            </Snackbar>
        </Box>
    )
}

export default NoteList
