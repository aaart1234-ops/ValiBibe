// features/note/components/SwipeableNoteRow.tsx
import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import {
    Grid,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
} from '@mui/material'
import {
    SwipeableList,
    SwipeableListItem,
    LeadingActions,
    TrailingActions,
    SwipeAction,
} from 'react-swipeable-list'
import 'react-swipeable-list/dist/styles.css'

import { Note, useDeleteNoteMutation } from '../noteApi'
import NoteRow from './NoteRow'

interface SwipeableNoteRowProps {
    note: Note
    onRefetch?: () => any
    onArchiveRequest?: (note: Note) => any
    onUnarchiveRequest?: (note: Note) => any
}

const SwipeableNoteRow: React.FC<SwipeableNoteRowProps> = ({
                                                               note,
                                                               onRefetch,
                                                               onArchiveRequest,
                                                               onUnarchiveRequest,
                                                           }) => {
    const [deleteNote] = useDeleteNoteMutation()
    const [confirmDialog, setConfirmDialog] = useState<'archive' | 'delete' | null>(null)
    const [isProcessing, setIsProcessing] = useState(false)

    const handleArchiveToggle = () => {
        if (note.archived) {
            onUnarchiveRequest?.(note)
        } else {
            // делегируем запрос архивирования родителю (чтобы он показал snackbar/undo)
            onArchiveRequest?.(note)
        }
    }

    const handleDeleteRequest = () => {
        setConfirmDialog('delete')
    }

    const handleConfirmDelete = async () => {
        if (isProcessing) return
        setIsProcessing(true)
        try {
            await deleteNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при удалении:', e)
        } finally {
            setIsProcessing(false)
            setConfirmDialog(null)
        }
    }

    return (
        <>
            <SwipeableList threshold={0.25}>
                <SwipeableListItem
                    leadingActions={
                        <LeadingActions>
                            <SwipeAction onClick={handleArchiveToggle}>
                                <Grid
                                    container
                                    key={note.id}
                                    sx={{
                                        width: { xs: '100%', sm: '48%', md: '32%' },
                                        backgroundColor: 'rgb(175, 238, 175)',
                                    }}
                                    display="flex"
                                    alignItems="center"
                                    justifyContent="center"
                                >
                                    {note.archived ? 'Извлечь из архива' : 'В архив'}
                                </Grid>
                            </SwipeAction>
                        </LeadingActions>
                    }
                    trailingActions={
                        <TrailingActions>
                            <SwipeAction onClick={handleDeleteRequest}>
                                <Grid
                                    container
                                    key={note.id}
                                    sx={{
                                        width: { xs: '100%', sm: '48%', md: '32%' },
                                        backgroundColor: 'rgb(238, 175, 175)',
                                    }}
                                    display="flex"
                                    alignItems="center"
                                    justifyContent="center"
                                >
                                    Удалить
                                </Grid>
                            </SwipeAction>
                        </TrailingActions>
                    }
                    blockSwipe={!!confirmDialog}
                >
                    <Link
                        to={`/notes/${note.id}`}
                        style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}
                    >
                        <NoteRow note={note} />
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
                        onClick={handleConfirmDelete}
                        color="error"
                        disabled={isProcessing}
                    >
                        Удалить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )
}

export default SwipeableNoteRow
