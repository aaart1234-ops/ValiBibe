import React from 'react'
import { Snackbar, Alert, Button } from '@mui/material'

type Props = {
    open: boolean
    autoHideDuration: number
    onClose: (event?: any, reason?: string) => void
    onUndo: () => void
    message?: string
}

const UndoSnackbar: React.FC<Props> = ({
                                           open,
                                           autoHideDuration,
                                           onClose,
                                           onUndo,
                                           message = 'Заметка перемещена в архив',
                                       }) => {
    return (
        <Snackbar
            open={open}
            autoHideDuration={autoHideDuration}
            onClose={onClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
            <Alert
                severity="info"
                sx={{ width: '100%' }}
                action={
                    <Button color="inherit" size="small" onClick={onUndo}>
                        Отменить
                    </Button>
                }
            >
                {message}
            </Alert>
        </Snackbar>
    )
}

export default UndoSnackbar
