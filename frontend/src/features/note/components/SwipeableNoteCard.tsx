// features/note/components/SwipeableNoteCard.tsx
import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
    Box,
} from '@mui/material'
import {
    SwipeableList,
    SwipeableListItem,
    LeadingActions,
    TrailingActions,
    SwipeAction,
} from 'react-swipeable-list'
import 'react-swipeable-list/dist/styles.css'

import { Note, useDeleteNoteMutation, useUnarchiveNoteMutation } from '../noteApi'
import NoteCard from './NoteCard'

interface SwipeableNoteCardProps {
    note: Note
    onRefetch?: () => void
    // новый callback: попросить внешний компонент (NoteList) начать процесс архивирования
    onRequestArchive?: (note: Note) => void
}

const SwipeableNoteCard: React.FC<SwipeableNoteCardProps> = ({ note, onRefetch, onRequestArchive }) => {
    const [deleteNote] = useDeleteNoteMutation()
    const [unarchiveNote] = useUnarchiveNoteMutation()

    const [confirmDelete, setConfirmDelete] = useState(false)
    const [isProcessing, setIsProcessing] = useState(false)

    const handleUnarchive = async () => {
        try {
            await unarchiveNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при разархивировании:', e)
        }
    }

    const handleDelete = async () => {
        if (isProcessing) return
        setIsProcessing(true)
        try {
            await deleteNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при удалении:', e)
        } finally {
            setConfirmDelete(false)
            setIsProcessing(false)
        }
    }

    return (
        <>
            <SwipeableList threshold={0.25}>
                <SwipeableListItem
                    leadingActions={
                        <LeadingActions>
                            <SwipeAction
                                onClick={() =>
                                    note.archived ? handleUnarchive() : onRequestArchive?.(note)
                                }
                            >
                                <Box
                                    display="flex"
                                    justifyContent="center"
                                    alignItems="center"
                                    bgcolor="rgb(175, 238, 175)"
                                    width="100%"
                                    height="100%"
                                >
                                    {note.archived ? 'Из архива' : 'В архив'}
                                </Box>
                            </SwipeAction>
                        </LeadingActions>
                    }
                    trailingActions={
                        <TrailingActions>
                            <SwipeAction onClick={() => setConfirmDelete(true)}>
                                <Box
                                    display="flex"
                                    justifyContent="center"
                                    alignItems="center"
                                    bgcolor="rgb(238, 175, 175)"
                                    width="100%"
                                    height="100%"
                                >
                                    Удалить
                                </Box>
                            </SwipeAction>
                        </TrailingActions>
                    }
                    blockSwipe={confirmDelete}
                >
                    <Link to={`/notes/${note.id}`} style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}>
                        <NoteCard note={note} />
                    </Link>
                </SwipeableListItem>
            </SwipeableList>

            {/* Диалог удаления — без изменений */}
            <Dialog open={confirmDelete} onClose={() => setConfirmDelete(false)}>
                <DialogTitle>Удалить заметку?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить эту заметку? Это действие необратимо.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDelete(false)} disabled={isProcessing}>
                        Отмена
                    </Button>
                    <Button onClick={handleDelete} color="error" disabled={isProcessing}>
                        Удалить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )
}

export default SwipeableNoteCard
