// NoteCard.tsx
import { Card, CardContent, Typography, Box, LinearProgress, Chip } from '@mui/material'
import { Note } from '../noteApi'
import dayjs from 'dayjs'
import EventIcon from '@mui/icons-material/Event'
import ArchiveIcon from '@mui/icons-material/Archive'

const NoteCard = ({ note }: { note: Note }) => (
    <Card sx={{ display: 'flex', flexDirection: 'column', flexGrow: 1, position: 'relative' }}>
        {/* Лейбл "В архиве" */}
        {note.archived && (
            <Chip
                icon={<ArchiveIcon />}
                label="В архиве"
                size="small"
                color="default"
                sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1 }}
                onClick={(e) => e.stopPropagation()}
            />
        )}
        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
            <Typography
                variant="h6"
                sx={{
                    display: '-webkit-box',
                    WebkitLineClamp: 2,
                    WebkitBoxOrient: 'vertical',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                }}
            >
                {note.title}
            </Typography>
            <Box display="flex" alignItems="center" gap={0.5} mt={0.5} mb={4}>
                <EventIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                <Typography variant="body2" color="text.secondary">
                    {dayjs(note.created_at).format('D MMMM YYYY')}
                </Typography>
            </Box>
            <Box mt="auto">
                <Typography variant="body2" color="text.secondary">
                    🧠 Уровень запоминания: {note.memoryLevel}%
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
