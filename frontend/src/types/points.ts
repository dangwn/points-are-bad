import { SessionUser } from './user'

export type SessionUserPoints = {
  points: number,
  correct_scores: number,
  largest_error: number,
  user: SessionUser
}

export type LeaderboardPoints = SessionUserPoints & {
  position: number|null
}