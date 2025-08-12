// components/NoteList.tsx
import React, { useState, useEffect, useMemo, useRef, useCallback } from 'react'
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
    Note,
} from '@/features/note/noteApi'

import NoteCard from '@/features/note/components/NoteCard'
import NoteRow from '@/features/note/components/NoteRow'
import SwipeableNoteCard from '@/features/note/components/SwipeableNoteCard'
import SwipeableNoteRow from '@/features/note/components/SwipeableNoteRow'
import NoteCreateDialog from '@/features/note/components/NoteCreateDialog'
import useMediaQuery from '@mui/material/useMediaQuery'
import { useTheme } from '@mui/material/styles'

/**
 * NoteList: хранит Snackbar/Undo/оптимистичное удаление для архивирования.
 *
 * UX:
 * - при запросе архивирования: заметка скрывается (optimistic), Snackbar появляется.
 * - если пользователь отменил — заметка возвращается (без обращения к серверу).
 * - если пользователь не отменил — по таймеру отправляем запрос archiveNote.
 * - если нет сети / запрос упал — сохраняем действие в localStorage (очередь) и
 *   пытаемся выполнить позже при событии "online".
 */

const ACTION_QUEUE_KEY = 'note_action_queue'
const SNACKBAR_TIMEOUT_MS = 4000

type PendingAction = {
    type: 'archive'
    noteId: string
    createdAt: number
}

const NoteList: React.FC = () => {
    const dispatch = useAppDispatch()
    const navigate = useNavigate()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { token } = useAppSelector((state) => state.auth)

    const [searchParams, setSearchParams] = useSearchParams()

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    // query params
    const showArchived = searchParams.get('archived') === 'true'
    const pageParam = parseInt(searchParams.get('page') || '1', 10)
    const limitParam = parseInt(searchParams.get('limit') || '10', 10)

    // фильтры / сортировка / поиск
    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
    const [searchQuery, setSearchQuery] = useState('')
    const [debouncedSearchQuery] = useDebounce(searchQuery, 300)

    const [limit, setLimit] = useState(limitParam)
    const [page, setPage] = useState(pageParam - 1)
    const offset = page * limit

    // optimistic архивация / Snackbar / undo
    const [optimisticArchivedIds, setOptimisticArchivedIds] = useState<Set<string>>(new Set())
    const undoTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const [snackbarOpen, setSnackbarOpen] = useState(false)
    const [lastArchivedNote, setLastArchivedNote] = useState<Note | null>(null)

    // RTK mutation для фактической архивации
    const [archiveNote] = useArchiveNoteMutation()

    // RTK Query получения списка
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
    }, [token]) // eslint-disable-line

    useEffect(() => {
        const params: Record<string, string> = {
            page: (page + 1).toString(),
            limit: limit.toString(),
        }
        if (showArchived) params.archived = 'true'
        setSearchParams(params, { replace: false })
    }, [page, limit, showArchived]) // eslint-disable-line

    // ---- localStorage queue helpers ----
    const getQueue = useCallback((): PendingAction[] => {
        try {
            return JSON.parse(localStorage.getItem(ACTION_QUEUE_KEY) || '[]') as PendingAction[]
        } catch {
            return []
        }
    }, [])

    const setQueue = useCallback((q: PendingAction[]) => {
        localStorage.setItem(ACTION_QUEUE_KEY, JSON.stringify(q))
    }, [])

    const pushToQueue = useCallback((action: PendingAction) => {
        const q = getQueue()
        q.push(action)
        setQueue(q)
    }, [getQueue, setQueue])

    // попытаемся выполнить очередь — вызывается при переходе в online
    const processQueue = useCallback(async () => {
        const q = getQueue()
        if (q.length === 0) return

        const remaining: PendingAction[] = []
        for (const action of q) {
            try {
                if (action.type === 'archive') {
                    await archiveNote(action.noteId).unwrap()
                }
                // если успешен — ничего не добавляем в remaining
            } catch (err) {
                // если сеть отсутствует — прекращаем попытки (оставим всё в очереди)
                if (!navigator.onLine) {
                    remaining.push(action, ...q.slice(q.indexOf(action) + 1))
                    break
                } else {
                    // если другая ошибка — оставляем действие, чтобы попробовать позже
                    remaining.push(action)
                }
            }
        }
        setQueue(remaining)
        if (remaining.length === 0) {
            // обновим UI
            refetch()
        }
    }, [archiveNote, getQueue, setQueue, refetch])

    useEffect(() => {
        // при появлении сети — пытаемся обработать очередь
        const onOnline = () => {
            processQueue().catch((e) => console.error('processQueue error', e))
        }
        window.addEventListener('online', onOnline)
        return () => window.removeEventListener('online', onOnline)
    }, [processQueue])

    // ---- archive request flow (from child Swipeable components) ----
    const handleRequestArchive = useCallback((note: Note) => {
        // пометим как оптимистично удалённую и покажем snackbar
        setOptimisticArchivedIds((prev) => {
            const next = new Set(prev)
            next.add(note.id)
            return next
        })
        setLastArchivedNote(note)
        setSnackbarOpen(true)

        // если уже был таймер — очистим
        if (undoTimerRef.current) {
            clearTimeout(undoTimerRef.current)
            undoTimerRef.current = null
        }

        // запланируем реальную архивацию через SNACKBAR_TIMEOUT_MS
        undoTimerRef.current = setTimeout(async () => {
            try {
                // попытаемся сразу архивировать
                await archiveNote(note.id).unwrap()
                // после успеха — обновим список
                refetch()
            } catch (err) {
                // если офлайн или ошибка — положим в очередь для поздней синхронизации
                pushToQueue({
                    type: 'archive',
                    noteId: note.id,
                    createdAt: Date.now(),
                })
            } finally {
                // очистим оптимистичные метки и последний note
                setOptimisticArchivedIds((prev) => {
                    const next = new Set(prev)
                    next.delete(note.id)
                    return next
                })
                setLastArchivedNote(null)
                setSnackbarOpen(false)
                undoTimerRef.current = null
            }
        }, SNACKBAR_TIMEOUT_MS)
    }, [archiveNote, refetch, pushToQueue])

    const handleUndo = useCallback(async () => {
        // отменяем запланированную архивацию
        if (undoTimerRef.current) {
            clearTimeout(undoTimerRef.current)
            undoTimerRef.current = null
        }
        // возвращаем заметку в UI
        setOptimisticArchivedIds((prev) => {
            const next = new Set(prev)
            if (lastArchivedNote) next.delete(lastArchivedNote.id)
            return next
        })
        setLastArchivedNote(null)
        setSnackbarOpen(false)
    }, [lastArchivedNote])

    const handleSnackbarClose = useCallback((_ev?: any, reason?: string) => {
        if (reason === 'clickaway') return
        setSnackbarOpen(false)
    }, [])

    // ---- rendering ----
    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box
            mt={4}
            pb={4}
            sx={{
                pl: isMobile ? 1 : 4,
                pr: isMobile ? 1 : 4,
            }}
        >
            <Box
                display="flex"
                flexDirection="column"
                justifyContent="space-between"
                alignItems="space-between"
                gap={2}
                mb={2}
            >
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

                        <Tooltip title={`Сортировать по ${sortDirection === 'asc' ? 'возрастанию' : 'убыванию'}`} sx={{ backgroundColor: '#e0e0e0' }} color="primary">
                            <IconButton onClick={() => setSortDirection((p) => (p === 'asc' ? 'desc' : 'asc'))}>
                                {sortDirection === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                            </IconButton>
                        </Tooltip>

                        <Tooltip title={viewMode === 'card' ? 'Список' : 'Карточки'} sx={{ backgroundColor: '#e0e0e0' }} color="primary">
                            <IconButton onClick={() => dispatch(toggleViewMode())}>
                                {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                            </IconButton>
                        </Tooltip>
                        {isMobile ? (
                            <Tooltip title={showArchived ? 'Показать активные' : 'Показать архив'} sx={{ backgroundColor: '#e0e0e0' }}>
                                <IconButton color="secondary" onClick={() => {
                                    setPage(0)
                                    setSearchParams((prev) => {
                                        const updated = new URLSearchParams(prev)
                                        if (showArchived) updated.delete('archived')
                                        else updated.set('archived', 'true')
                                        updated.set('page', '1')
                                        return updated
                                    })
                                }}>
                                    {showArchived ? <ArchiveOutlinedIcon /> : <ArchiveIcon />}
                                </IconButton>
                            </Tooltip>
                        ) : (
                            <Button variant={showArchived ? 'contained' : 'outlined'} color="primary" onClick={() => {
                                setPage(0)
                                setSearchParams((prev) => {
                                    const updated = new URLSearchParams(prev)
                                    if (showArchived) updated.delete('archived')
                                    else updated.set('archived', 'true')
                                    updated.set('page', '1')
                                    return updated
                                })
                            }}>
                                {showArchived ? 'Показать активные' : 'Показать архив'}
                            </Button>
                        )}
                    </Box>

                    {!isMobile && (
                        <IconButton onClick={() => setOpenCreateDialog(true)} color="primary" size="large" sx={{ padding: 0, '&:hover': { backgroundColor: '#e0e0e0' } }}>
                            <AddIcon fontSize="large" />
                        </IconButton>
                    )}
                </Box>


                {/* FAB на мобильных */}
                {isMobile && (
                    <Fab color="primary" onClick={() => setOpenCreateDialog(true)} sx={{ position: 'fixed', bottom: 24, right: 24, zIndex: 10 }}>
                        <AddIcon />
                    </Fab>
                )}

            </Box>

            {notes.length === 0 ? (
                <Box textAlign="center" mt={10}>
                    <Typography variant="h6" gutterBottom>У вас пока нет заметок</Typography>
                    <Typography variant="body2" color="text.secondary" mb={2}>Создайте первую заметку, чтобы начать тренироваться</Typography>
                    <Button variant="contained" color="primary" component={Link} to="/notes/create">Создать заметку</Button>
                </Box>
            ) : viewMode === 'card' ? (
                <Grid container spacing={2} justifyContent="start" mt={4}>
                    {notes
                        .filter((note) => !optimisticArchivedIds.has(note.id))
                        .map((note) => (
                            <Grid key={note.id} sx={{ width: { xs: '100%', sm: '48%', md: '32%' } }} display="flex">
                                {isMobile ? (
                                    <SwipeableNoteCard note={note} onRefetch={refetch} onRequestArchive={handleRequestArchive} />
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
                    {notes
                        .filter((note) => !optimisticArchivedIds.has(note.id))
                        .map((note) =>
                            isMobile ? (
                                <SwipeableNoteRow key={note.id} note={note} onRefetch={refetch} onRequestArchive={handleRequestArchive} />
                            ) : (
                                <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                                    <NoteRow note={note} />
                                </Link>
                            )
                        )}
                </Box>
            )}

            {/* Pagination */}
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

            <NoteCreateDialog open={openCreateDialog} onClose={() => setOpenCreateDialog(false)} onCreated={refetch} />

            {/* Snackbar для архивации с Undo */}
            <Snackbar
                open={snackbarOpen}
                autoHideDuration={SNACKBAR_TIMEOUT_MS}
                onClose={handleSnackbarClose}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert
                    severity="info"
                    sx={{ width: '100%' }}
                    action={
                        <Button color="inherit" size="small" onClick={handleUndo}>
                            Отменить
                        </Button>
                    }
                >
                    Заметка перемещена в архив
                </Alert>
            </Snackbar>
        </Box>
    )
}

export default NoteList
