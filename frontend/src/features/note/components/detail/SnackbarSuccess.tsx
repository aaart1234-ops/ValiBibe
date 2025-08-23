import { Snackbar, Alert, Button } from '@mui/material'

interface SnackbarSuccessProps {
    open: boolean
    message: string
    onClose: () => void
    actionText?: string
    onAction?: () => void
}

export const SnackbarSuccess = ({
                                    open,
                                    message,
                                    onClose,
                                    actionText,
                                    onAction,
                                }: SnackbarSuccessProps) => {
    return (
        <Snackbar open={open} autoHideDuration={5000} onClose={onClose}>
            <Alert
                severity="success"
                action={
                    actionText ? (
                        <Button color="inherit" size="small" onClick={onAction}>
                            {actionText}
                        </Button>
                    ) : null
                }
            >
                {message}
            </Alert>
        </Snackbar>
    )
}
