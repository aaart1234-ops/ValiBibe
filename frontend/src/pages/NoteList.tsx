// components/NoteList.tsx
import { useAppDispatch, useAppSelector } from '@/app/hooks'
import { toggleViewMode } from '@/features/note/noteSlice'
import { useGetNotesQuery } from '@/features/note/noteApi'
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
    Tooltip
} from '@mui/material'
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward'
import Grid from '@mui/material/Grid'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import NoteCard from '../features/note/components/NoteCard'
import NoteRow from '../features/note/components/NoteRow'
import React, {useState, useMemo, useEffect} from "react"
import { Link } from 'react-router-dom'
import AddIcon from '@mui/icons-material/Add'
import NoteCreateDialog from '@/features/note/components/NoteCreateDialog'
import { Fab } from '@mui/material'




const NoteList = () => {
    const dispatch = useAppDispatch()
    const viewMode = useAppSelector((state) => state.notes.viewMode)
    const { data: notes, isLoading, isError, refetch } = useGetNotesQuery(undefined, {
        refetchOnMountOrArgChange: true, // полезно в любом случае
    })

    const { token } = useAppSelector(state => state.auth)

    const [openCreateDialog, setOpenCreateDialog] = useState(false)

    useEffect(() => {
        if (token) {
            refetch()
        }
    }, [token])

    const [sortBy, setSortBy] = useState<'created_at' | 'next_review_at'>('created_at')

    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc') // новое состояние
    const toggleSortDirection = () => {
        setSortDirection(prev => (prev === 'asc' ? 'desc' : 'asc'))
    }

    const [searchQuery, setSearchQuery] = useState('')

    const filteredAndSortedNotes = useMemo(() => {
        if (!notes) return []

        return [...notes]
            .filter(note =>
                note.title.toLowerCase().includes(searchQuery.toLowerCase())
            )
            .sort((a, b) => {
                const fieldA = new Date(a[sortBy]!).getTime()
                const fieldB = new Date(b[sortBy]!).getTime()
                return sortDirection === 'asc' ? fieldA - fieldB : fieldB - fieldA
            })
    }, [notes, sortBy, searchQuery, sortDirection])

    if (isLoading) return <CircularProgress />
    if (isError || !notes) return <Typography>Ошибка загрузки заметок</Typography>

    return (
        <Box mt={4}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2} flexWrap="wrap" gap={2}>
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
                            }
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

            {filteredAndSortedNotes.length === 0 ? (
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
                <Grid container spacing={1} justifyContent="space-between" mt={8}>
                    {filteredAndSortedNotes.map(note => (
                        <Grid
                            key={note.id}
                            sx={{ width: { xs: '100%', sm: '31%', md: '31%' } }}
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
                    {filteredAndSortedNotes.map(note => (
                        <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                            <NoteRow note={note} />
                        </Link>
                    ))}
                </Box>
            )}
            <NoteCreateDialog
                open={openCreateDialog}
                onClose={() => setOpenCreateDialog(false)}
                onCreated={refetch}
            />
        </Box>
    )
}

export default NoteList
