import 'dotenv/config'
import express from 'express'
import cors from 'cors'
import { aiRouter } from './routes/ai.js'

const app = express()
const PORT = parseInt(process.env.PORT || '3001', 10)

app.use(cors())
app.use(express.json())

app.use('/api/ai', aiRouter)

app.get('/api/health', (_req, res) => {
  res.json({ status: 'ok', service: 'ai-platform-backend' })
})

app.listen(PORT, () => {
  console.log(`Backend server running on http://localhost:${PORT}`)
})
