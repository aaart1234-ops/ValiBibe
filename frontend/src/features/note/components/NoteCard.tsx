// NoteCard.tsx
import { Card, CardContent, Typography, Box, LinearProgress } from '@mui/material'
import { Note } from '../noteApi'
import dayjs from 'dayjs'

const NoteCard = ({ note }: { note: Note }) => (
    <Card sx={{ display: 'flex', flexDirection: 'column', flexGrow: 1 }}>
        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6">{note.title}</Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
                Создана: {dayjs(note.created_at).format('D MMMM YYYY')}
            </Typography>
            <Box mt="auto">
                <Typography variant="body2" color="text.secondary">
                    Уровень запоминания: {note.memoryLevel}%
                </Typography>
                <LinearProgress
                    variant="determinate"
                    value={note.memoryLevel}
                    sx={{ height: 8, borderRadius: 5, mt: 0.5 }}
                    color={
                        note.memoryLevel < 40
                            ? 'error'
                            : note.memoryLevel < 70
                                ? 'warning'
                                : 'success'
                    }
                />
            </Box>
        </CardContent>
    </Card>
)

export default NoteCard
