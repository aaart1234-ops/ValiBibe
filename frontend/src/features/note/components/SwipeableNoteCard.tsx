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
    SwipeAction
} from 'react-swipeable-list'
import 'react-swipeable-list/dist/styles.css'

import { Note, useArchiveNoteMutation, useDeleteNoteMutation, useUnarchiveNoteMutation } from '../noteApi'
import NoteCard from './NoteCard'

interface SwipeableNoteCardProps {
    note: Note
    onRefetch?: () => void
}

const SwipeableNoteCard: React.FC<SwipeableNoteCardProps> = ({ note, onRefetch }) => {
    const [archiveNote] = useArchiveNoteMutation()
    const [unarchiveNote] = useUnarchiveNoteMutation()
    const [deleteNote] = useDeleteNoteMutation()

    const [confirmDialog, setConfirmDialog] = useState<null | 'delete' | 'archive'>(null)
    const [isProcessing, setIsProcessing] = useState(false)

    const handleConfirm = async () => {
        if (isProcessing) return
        setIsProcessing(true)
        try {
            if (confirmDialog === 'delete') {
                await deleteNote(note.id).unwrap()
            } else if (confirmDialog === 'archive') {
                await archiveNote(note.id).unwrap()
            }
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при подтверждении действия:', e)
        } finally {
            setConfirmDialog(null)
            setIsProcessing(false)
        }
    }

    const handleUnarchive = async () => {
        try {
            await unarchiveNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при разархивировании:', e)
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
                                    note.archived
                                        ? handleUnarchive()
                                        : setConfirmDialog('archive')
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
                            <SwipeAction
                                onClick={() => setConfirmDialog('delete')}
                            >
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
                    blockSwipe={!!confirmDialog}
                >
                    <Link
                        to={`/notes/${note.id}`}
                        style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}
                    >
                        <NoteCard note={note} />
                    </Link>
                </SwipeableListItem>
            </SwipeableList>

            <Dialog open={!!confirmDialog} onClose={() => setConfirmDialog(null)}>
                <DialogTitle>
                    {confirmDialog === 'delete' ? 'Удалить заметку?' : 'Архивировать заметку?'}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmDialog === 'delete'
                            ? 'Вы уверены, что хотите удалить эту заметку? Это действие необратимо.'
                            : 'После архивирования заметка будет скрыта из основного списка.'}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDialog(null)} disabled={isProcessing}>
                        Отмена
                    </Button>
                    <Button
                        onClick={handleConfirm}
                        color={confirmDialog === 'delete' ? 'error' : 'primary'}
                        disabled={isProcessing}
                    >
                        Подтвердить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )
}

export default SwipeableNoteCard
