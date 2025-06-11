// NoteRow.tsx
import { Box, Typography, Divider } from '@mui/material'
import { Note } from '../noteApi'
import dayjs from 'dayjs'

const NoteRow = ({ note }: { note: Note }) => (
    <Box mb={1}>
        <Typography variant="subtitle1">{note.title}</Typography>
        <Typography variant="caption" color="text.secondary">
            {dayjs(note.created_at).format('D MMMM YYYY')}
        </Typography>
        <Divider sx={{ mt: 1 }} />
    </Box>
)

export default NoteRow
