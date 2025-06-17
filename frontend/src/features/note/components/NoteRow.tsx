import { Box, Typography, Divider, LinearProgress, Stack } from '@mui/material'
import { Note } from '../noteApi'
import dayjs from 'dayjs'

const NoteRow = ({ note }: { note: Note }) => (
    <Box mb={2}>
        <Stack direction="row" justifyContent="space-between" alignItems="center">
            <Typography variant="subtitle1" fontWeight={500}>
                {note.title}
            </Typography>
            <Typography variant="caption" color="text.secondary">
                {dayjs(note.created_at).format('D MMM YYYY')}
            </Typography>
        </Stack>

        <Box mt={1}>
            <Typography variant="caption" color="text.secondary">
                Уровень запоминания: {note.memoryLevel}%
            </Typography>
            <LinearProgress
                variant="determinate"
                value={note.memoryLevel}
                sx={{ height: 6, borderRadius: 4, mt: 0.5 }}
                color={note.memoryLevel < 40 ? 'error' : note.memoryLevel < 70 ? 'warning' : 'success'}
            />
        </Box>

        <Divider sx={{ mt: 2 }} />
    </Box>
)

export default NoteRow
