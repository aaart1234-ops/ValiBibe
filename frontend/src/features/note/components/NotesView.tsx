import React from 'react'
import { Link } from 'react-router-dom'
import { Box, Grid, Typography } from '@mui/material'
import { Note } from '@/features/note/noteApi'
import NoteCard from '@/features/note/components/NoteCard'
import NoteRow from '@/features/note/components/NoteRow'
import SwipeableNoteCard from '@/features/note/components/SwipeableNoteCard'
import SwipeableNoteRow from '@/features/note/components/SwipeableNoteRow'

type Props = {
    notes: Note[]
    viewMode: 'card' | 'list'
    isMobile: boolean
    onRequestArchive: (note: Note) => void
    onRefetch: () => void
    isArchiveView: boolean
}

const NotesView: React.FC<Props> = ({
                                        notes,
                                        viewMode,
                                        isMobile,
                                        onRequestArchive,
                                        onRefetch,
                                        isArchiveView
                                    }) => {
    if (notes.length === 0) {
        return (
            <Box textAlign="center" mt={10}>
                <Typography variant="h6" gutterBottom>
                    {isArchiveView ? "В архиве нет заметок" : "У вас пока нет заметок"}
                </Typography>
                <Typography variant="body2" color="text.secondary" mb={2}>
                    Создайте первую заметку, чтобы начать тренироваться
                </Typography>
                <Link to="/notes/create" style={{ textDecoration: 'none' }}>
                    <Typography
                        component="span"
                        sx={{
                            display: 'inline-block',
                            px: 2, py: 1,
                            bgcolor: 'primary.main',
                            color: 'primary.contrastText',
                            borderRadius: 1
                        }}
                    >
                        Создать заметку
                    </Typography>
                </Link>
            </Box>
        )
    }

    if (viewMode === 'card') {
        return (
            <Grid container spacing={2} justifyContent="start" mt={4}>
                {notes.map((note) => (
                    <Grid key={note.id} sx={{ width: { xs: '100%', sm: '48%', md: '32%' } }} display="flex">
                        {isMobile ? (
                            <SwipeableNoteCard note={note} onRefetch={onRefetch} onRequestArchive={onRequestArchive} />
                        ) : (
                            <Link to={`/notes/${note.id}`} style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}>
                                <NoteCard note={note} />
                            </Link>
                        )}
                    </Grid>
                ))}
            </Grid>
        )
    }

    // list
    return (
        <Box mt={4}>
            {notes.map((note) =>
                isMobile ? (
                    <SwipeableNoteRow key={note.id} note={note} onRefetch={onRefetch} onRequestArchive={onRequestArchive} />
                ) : (
                    <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none' }}>
                        <NoteRow note={note} />
                    </Link>
                )
            )}
        </Box>
    )
}

export default NotesView
