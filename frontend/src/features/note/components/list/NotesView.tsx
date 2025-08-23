import React from 'react'
import { Box, Grid, Typography } from '@mui/material'
import { Note } from '@/features/note/noteApi'
import NoteCard from '@/features/note/components/list/NoteCard'
import NoteRow from '@/features/note/components/list/NoteRow'
import SwipeableNoteCard from '@/features/note/components/list/SwipeableNoteCard'
import SwipeableNoteRow from '@/features/note/components/list/SwipeableNoteRow'
import ArchiveOutlined from '@mui/icons-material/ArchiveOutlined'
import NoteOutlined from '@mui/icons-material/NoteOutlined'

type Props = {
    notes: Note[]
    viewMode: 'card' | 'list'
    isMobile: boolean
    onRequestArchive: (note: Note) => void
    onRequestDelete: (note: Note) => void
    onRefetch: () => void
    isArchiveView: boolean
}

const NotesView: React.FC<Props> = ({
                                        notes,
                                        viewMode,
                                        isMobile,
                                        onRequestArchive,
                                        onRequestDelete,
                                        onRefetch,
                                        isArchiveView,
                                    }) => {
    if (notes.length === 0) {
        return (
            <Box
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    py: 6,
                    opacity: 0.7,
                }}
            >
                {isArchiveView ? (
                    <>
                        <ArchiveOutlined sx={{ fontSize: 60, mb: 2 }} color="action" />
                        <Typography variant="h6">В архиве нет заметок</Typography>
                    </>
                ) : (
                    <>
                        <NoteOutlined sx={{ fontSize: 60, mb: 2 }} color="action" />
                        <Typography variant="h6">У вас пока нет заметок</Typography>
                        <Typography variant="body2" color="text.secondary">
                            Нажмите «+» чтобы создать первую
                        </Typography>
                    </>
                )}
            </Box>
        )
    }

    if (viewMode === 'card') {
        return (
            <Grid container spacing={2} justifyContent="start" mt={4}>
                {notes.map((note) => (
                    <Grid
                        key={note.id}
                        sx={{ width: { xs: '100%', sm: '48%', md: '32%' } }}
                        display="flex"
                    >
                        {isMobile ? (
                            <SwipeableNoteCard
                                note={note}
                                onRefetch={onRefetch}
                                onRequestArchive={onRequestArchive}
                            />
                        ) : (
                            <NoteCard
                                note={note}
                                isArchiveView={isArchiveView}
                                onRequestArchive={onRequestArchive}
                            />
                        )}
                    </Grid>
                ))}
            </Grid>
        )
    }

    // list
    return (
        <Box mt={4} sx={{ pl: isMobile ? 1 : 4, pr: isMobile ? 1 : 4 }}>
            {notes.map((note) =>
                isMobile ? (
                    <SwipeableNoteRow
                        key={note.id}
                        note={note}
                        onRefetch={onRefetch}
                        onRequestArchive={onRequestArchive}
                        onRequestDelete={onRequestDelete}
                    />
                ) : (
                    <NoteRow
                        key={note.id}
                        note={note}
                        isArchiveView={isArchiveView}
                        onRequestArchive={onRequestArchive}
                        onRequestDelete={onRequestDelete}
                    />
                )
            )}
        </Box>
    )
}

export default NotesView
