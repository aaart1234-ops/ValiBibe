import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
} from '@mui/material'

interface ConfirmDialogProps {
    type: 'delete' | 'archive' | null
    open: boolean
    onClose: () => void
    onConfirm: () => void
}

export const ConfirmDialog = ({ type, open, onClose, onConfirm }: ConfirmDialogProps) => {
    if (!type) return null
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>{type === 'delete' ? 'Удалить заметку?' : 'Архивировать заметку?'}</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    {type === 'delete'
                        ? 'Вы уверены, что хотите удалить эту заметку? Это действие необратимо.'
                        : 'После архивирования заметка будет скрыта из основного списка.'}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Отмена</Button>
                <Button onClick={onConfirm} color={type === 'delete' ? 'error' : 'primary'}>
                    Подтвердить
                </Button>
            </DialogActions>
        </Dialog>
    )
}
