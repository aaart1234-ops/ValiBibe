// features/note/components/NoteDeleteDialog.tsx
import React, { useState } from 'react'
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
} from '@mui/material'
import { Note, useDeleteNoteMutation } from '../../noteApi'

interface NoteDeleteDialogProps {
    note: Note | null
    open: boolean
    onClose: () => void
    onDeleted?: () => void
}

const NoteDeleteDialog: React.FC<NoteDeleteDialogProps> = ({ note, open, onClose, onDeleted }) => {
    const [deleteNote] = useDeleteNoteMutation()
    const [isProcessing, setIsProcessing] = useState(false)

    const handleDelete = async () => {
        if (!note || isProcessing) return
        setIsProcessing(true)
        try {
            await deleteNote(note.id).unwrap()
            onDeleted?.()
            onClose()
        } catch (e) {
            console.error('Ошибка при удалении:', e)
        } finally {
            setIsProcessing(false)
        }
    }

    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Удалить заметку?</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    Вы уверены, что хотите удалить эту заметку? Это действие необратимо.
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} disabled={isProcessing}>
                    Отмена
                </Button>
                <Button onClick={handleDelete} color="error" disabled={isProcessing}>
                    Удалить
                </Button>
            </DialogActions>
        </Dialog>
    )
}

export default NoteDeleteDialog
